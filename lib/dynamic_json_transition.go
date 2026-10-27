package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type deprecatedJSONEncodingValidator struct {
	rule, example, deadline string
	valueTypes              map[string]string
	requiredPaths           []string
	nestedPaths             []string
}

// DeprecatedJSONEncoding warns about JSON-encoded values. valueTypes maps fields of the typed schema to "integer",
// "boolean" or "list", with "[]" for list items and "*" for map values, for example {"user_ids[]": "integer"}. A
// configured value of another type gets guidance for its field, in the same warning when the value is also encoded.
// requiredPaths name fields, in the same form, that the typed schema requires although the API also accepts null for
// them. A configured null there gets guidance, because Terraform treats it as missing. An unknown value does not.
// nestedPaths name the discriminator of list union items that the typed schema moves under their variant name, for
// example "folder_behaviors[].behavior". An item that still has it uses the old nesting and gets example as guidance.
func DeprecatedJSONEncoding(rule, example, deadline string, valueTypes map[string]string, requiredPaths []string, nestedPaths ...string) validator.Dynamic {
	return deprecatedJSONEncodingValidator{rule: rule, example: example, deadline: deadline, valueTypes: valueTypes, requiredPaths: requiredPaths, nestedPaths: nestedPaths}
}

func (v deprecatedJSONEncodingValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v deprecatedJSONEncodingValidator) MarkdownDescription(_ context.Context) string {
	return "Use native HCL instead of a JSON-encoded string. JSON strings stop being accepted on " + v.deadline + "."
}

func (v deprecatedJSONEncodingValidator) ValidateDynamic(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.IsUnderlyingValueNull() || req.ConfigValue.IsUnderlyingValueUnknown() {
		return
	}
	_, encoded := req.ConfigValue.UnderlyingValue().(types.String)
	if !encoded && len(v.valueTypes) == 0 && len(v.requiredPaths) == 0 && len(v.nestedPaths) == 0 {
		return
	}
	value, diags := JSONValueToAPI(ctx, req.Path, req.ConfigValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	property := v.rule[strings.LastIndex(v.rule, ".")+1:]
	guidance := valueTypeGuidance(property, value, v.valueTypes)
	// JSONValueToAPI reads an unknown value as null, so null fields are found in the configuration itself. A decoded
	// JSON string has no unknown values.
	configured := req.ConfigValue
	if encoded && len(v.requiredPaths) > 0 {
		configured, diags = ToDynamic(ctx, req.Path, value, nil)
		resp.Diagnostics.Append(diags...)
	}
	nullGuidance := requiredNullGuidance(property, configured, v.requiredPaths, v.valueTypes, v.deadline)
	nesting := ""
	for _, nestedPath := range v.nestedPaths {
		if JSONValueHas(value, nestedPath) {
			field, discriminator, _ := strings.Cut(nestedPath, "[].")
			nesting = fmt.Sprintf("Move each %s item under the name in its %s field, as in this example: %s. Both nestings work now. The old nesting stops working on %s, when the typed schema ships in a major provider release.", field, discriminator, v.example, v.deadline)
		}
	}
	switch {
	case encoded && nesting != "":
		resp.Diagnostics.AddAttributeWarning(req.Path, "Deprecated JSON encoding: "+v.rule, joinSentences("Remove jsonencode(...) or replace the JSON string with native HCL. JSON strings stop being accepted on "+v.deadline+".", nesting, guidance, nullGuidance))
	case encoded:
		detail := fmt.Sprintf("Remove jsonencode(...) or replace the JSON string with native HCL. Keep the same keys and nesting. JSON strings stop being accepted on %s, when the typed schema ships in a major provider release. No additional wrapper is needed.", v.deadline)
		resp.Diagnostics.AddAttributeWarning(req.Path, "Deprecated JSON encoding: "+v.rule, joinSentences(detail, guidance, nullGuidance))
	case nesting != "":
		resp.Diagnostics.AddAttributeWarning(req.Path, "Deprecated nesting: "+v.rule, joinSentences(nesting, guidance, nullGuidance))
	case guidance != "":
		resp.Diagnostics.AddAttributeWarning(req.Path, "Typed values expected: "+v.rule, joinSentences(guidance, "The typed schema that ships on "+v.deadline+" uses these types.", nullGuidance))
	case nullGuidance != "":
		resp.Diagnostics.AddAttributeWarning(req.Path, "Typed values expected: "+v.rule, nullGuidance)
	}
}

// joinSentences joins the non-empty parts of a warning detail with single spaces.
func joinSentences(parts ...string) string {
	return strings.Join(slices.DeleteFunc(parts, func(part string) bool { return part == "" }), " ")
}

// requiredNullGuidance names the requiredPaths fields that the configuration sets to a known null.
func requiredNullGuidance(property string, configured attr.Value, requiredPaths []string, valueTypes map[string]string, deadline string) string {
	var fields []string
	replacement := "true or false"
	for _, requiredPath := range requiredPaths {
		if knownNullAt(configured, valuePathSegments(requiredPath)) {
			fields = append(fields, guidanceField(property, requiredPath))
			if valueTypes[requiredPath] != "boolean" {
				replacement = "values"
			}
		}
	}
	if len(fields) == 0 {
		return ""
	}
	return fmt.Sprintf("Replace null in %s with %s. The typed schema that ships on %s requires these fields and treats null as missing.", fieldList(fields), replacement, deadline)
}

// knownNullAt reports whether a value at the path segments is configured as null. A missing key or an unknown value is
// not null.
func knownNullAt(value attr.Value, segments []string) bool {
	if dynamic, ok := value.(types.Dynamic); ok && !dynamic.IsNull() && !dynamic.IsUnknown() {
		value = dynamic.UnderlyingValue()
	}
	if value == nil || value.IsUnknown() {
		return false
	}
	if value.IsNull() || len(segments) == 0 {
		return value.IsNull() && len(segments) == 0
	}
	var entries map[string]attr.Value
	var items []attr.Value
	switch value := value.(type) {
	case types.Object:
		entries = value.Attributes()
	case types.Map:
		entries = value.Elements()
	case types.Tuple:
		items = value.Elements()
	case types.List:
		items = value.Elements()
	case types.Set:
		items = value.Elements()
	}
	var children []attr.Value
	switch segments[0] {
	case "":
		children = []attr.Value{value}
	case "[]":
		children = items
	case "*":
		children = slices.Collect(maps.Values(entries))
	default:
		if entry, exists := entries[segments[0]]; exists {
			children = []attr.Value{entry}
		}
	}
	return slices.ContainsFunc(children, func(child attr.Value) bool { return knownNullAt(child, segments[1:]) })
}

type validJSONEncodingValidator struct{}

// ValidJSONEncoding rejects a JSON-encoded string that does not decode, at plan time, for a property that stays Dynamic
// in every release. Native values and valid JSON strings pass without a warning.
func ValidJSONEncoding() validator.Dynamic {
	return validJSONEncodingValidator{}
}

func (v validJSONEncodingValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v validJSONEncodingValidator) MarkdownDescription(_ context.Context) string {
	return "A JSON-encoded string must be valid JSON."
}

func (v validJSONEncodingValidator) ValidateDynamic(_ context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.IsUnderlyingValueNull() || req.ConfigValue.IsUnderlyingValueUnknown() {
		return
	}
	if encoded, ok := req.ConfigValue.UnderlyingValue().(types.String); ok && !json.Valid([]byte(encoded.ValueString())) {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid JSON value", "Expected valid JSON or a native HCL value. The JSON string could not be decoded.")
	}
}

type rejectJSONNullValidator struct{ detail string }

type rejectEmptyListValidator struct{ detail string }

// RejectEmptyList rejects an empty list for a property whose empty list the API ignores, so the value read back never
// matches it. It validates a typed list and a Dynamic value, native or JSON-encoded.
func RejectEmptyList(detail string) rejectEmptyListValidator {
	return rejectEmptyListValidator{detail: detail}
}

func (v rejectEmptyListValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v rejectEmptyListValidator) MarkdownDescription(_ context.Context) string {
	return "An empty list is not supported."
}

func (v rejectEmptyListValidator) ValidateList(_ context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() && len(req.ConfigValue.Elements()) == 0 {
		resp.Diagnostics.AddAttributeError(req.Path, "Unsupported empty list", v.detail)
	}
}

func (v rejectEmptyListValidator) ValidateDynamic(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.IsUnderlyingValueNull() || req.ConfigValue.IsUnderlyingValueUnknown() {
		return
	}
	value, diags := JSONValueToAPI(ctx, req.Path, req.ConfigValue)
	if list, ok := value.([]any); ok && !diags.HasError() && len(list) == 0 {
		resp.Diagnostics.AddAttributeError(req.Path, "Unsupported empty list", v.detail)
	}
}

type knownKeysValidator struct {
	detail string
	names  []string
}

// KnownKeys rejects a value that is not an object, and top-level keys outside names, for a Dynamic object whose unknown
// keys the API removes when it saves, so the value read back never matches the configuration.
func KnownKeys(detail string, names ...string) validator.Dynamic {
	return knownKeysValidator{detail: detail, names: names}
}

func (v knownKeysValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v knownKeysValidator) MarkdownDescription(_ context.Context) string {
	return "Every top-level key must be a documented field."
}

func (v knownKeysValidator) ValidateDynamic(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.IsUnderlyingValueNull() || req.ConfigValue.IsUnderlyingValueUnknown() {
		return
	}
	value, diags := JSONValueToAPI(ctx, req.Path, req.ConfigValue)
	// A JSON-encoded null is a supported value: Rails stores it as {}.
	if diags.HasError() || value == nil {
		return
	}
	object, ok := value.(map[string]any)
	if !ok {
		resp.Diagnostics.AddAttributeError(req.Path, "Unsupported value", fmt.Sprintf("%s must be an object of documented fields.", req.Path))
		return
	}
	var unknown []string
	for name := range object {
		if !slices.Contains(v.names, name) {
			unknown = append(unknown, name)
		}
	}
	if len(unknown) > 0 {
		slices.Sort(unknown)
		resp.Diagnostics.AddAttributeError(req.Path, "Unsupported keys", fmt.Sprintf("%s sets %s, which %s not documented. %s", req.Path, strings.Join(unknown, ", "), map[bool]string{true: "is", false: "are"}[len(unknown) == 1], v.detail))
	}
}

// RejectJSONNull rejects a JSON-encoded null, such as jsonencode(null), for a property whose null the API ignores, so
// the value read back never matches it. Native null and an omitted value are not affected.
func RejectJSONNull(detail string) validator.Dynamic {
	return rejectJSONNullValidator{detail: detail}
}

func (v rejectJSONNullValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v rejectJSONNullValidator) MarkdownDescription(_ context.Context) string {
	return "A JSON-encoded null is not supported."
}

func (v rejectJSONNullValidator) ValidateDynamic(_ context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.IsUnderlyingValueNull() || req.ConfigValue.IsUnderlyingValueUnknown() {
		return
	}
	if encoded, ok := req.ConfigValue.UnderlyingValue().(types.String); ok && strings.TrimSpace(encoded.ValueString()) == "null" {
		resp.Diagnostics.AddAttributeError(req.Path, "Unsupported JSON null", v.detail)
	}
}

// JSONValueHas reports whether value has a value other than null at valuePath, with "[]" for list items and "*" for
// map values, for example "folder_behaviors[].behavior".
func JSONValueHas(value any, valuePath string) bool {
	return !hasValueType(value, valuePathSegments(valuePath), "")
}

// valuePathSegments splits a value path at dots and before each "[]" or "[name=value]" item selector.
func valuePathSegments(valuePath string) []string {
	return strings.Split(strings.ReplaceAll(valuePath, "[", ".["), ".")
}

var valueTypeForms = []struct{ valueType, form string }{{"integer", "integers"}, {"boolean", "true or false"}, {"list", "lists"}, {"object", "objects"}}

var itemSelectors = regexp.MustCompile(`\[[^\]]*\]`)

// valueTypeGuidance names the fields whose configured value has another type than valueTypes gives them, with an
// example for the first field of each non-boolean type. It never repeats configured values or map keys.
func valueTypeGuidance(property string, value any, valueTypes map[string]string) string {
	var sentences []string
	for _, form := range valueTypeForms {
		var paths, fields []string
		for _, valuePath := range slices.Sorted(maps.Keys(valueTypes)) {
			if valueTypes[valuePath] == form.valueType && !hasValueType(value, valuePathSegments(valuePath), form.valueType) {
				field := guidanceField(property, valuePath)
				if !slices.Contains(fields, field) {
					paths = append(paths, valuePath)
					fields = append(fields, field)
				}
			}
		}
		if len(paths) == 0 {
			continue
		}
		sentence := "Write " + fieldList(fields) + " as " + form.form
		name := fields[0][strings.LastIndex(fields[0], ".")+1:]
		if example := valueExample(strings.TrimRight(paths[0], "[]"), valueTypes); example != "" && name != "*" && form.valueType != "boolean" {
			sentence += ", for example " + name + " = " + example
		}
		sentences = append(sentences, sentence+".")
	}
	return strings.Join(sentences, " ")
}

// guidanceField names the field at valuePath without item selectors. A map value is named after the property.
func guidanceField(property, valuePath string) string {
	field := strings.Trim(itemSelectors.ReplaceAllString(valuePath, ""), ".")
	if field == "" || strings.HasPrefix(field, "*") {
		field = strings.Trim(property+"."+field, ".")
	}
	return field
}

func fieldList(fields []string) string {
	if len(fields) == 1 {
		return fields[0]
	}
	return strings.Join(fields[:len(fields)-1], ", ") + " and " + fields[len(fields)-1]
}

// hasValueType reports whether every value at the path segments has valueType. Missing, null and unknown values, and
// shapes that do not match the path, are left to the API's validation.
func hasValueType(value any, segments []string, valueType string) bool {
	if len(segments) == 0 {
		switch value := value.(type) {
		case nil:
			return true
		case float64:
			return valueType == "integer" && value == math.Trunc(value)
		case int64:
			return valueType == "integer"
		case bool:
			return valueType == "boolean"
		case []any:
			return valueType == "list"
		case map[string]any:
			return valueType == "object"
		}
		return false
	}
	var values []any
	object, isObject := value.(map[string]any)
	switch segments[0] {
	case "":
		values = []any{value}
	case "[]":
		values, _ = value.([]any)
	case "*":
		values = slices.Collect(maps.Values(object))
	default:
		if selector, isSelector := strings.CutPrefix(segments[0], "["); isSelector {
			// "[name=value]" selects the list items, or the object, whose name field equals value.
			name, selected, _ := strings.Cut(strings.TrimSuffix(selector, "]"), "=")
			candidates, isList := value.([]any)
			if !isList {
				candidates = []any{value}
			}
			for _, candidate := range candidates {
				if item, isItem := candidate.(map[string]any); isItem && item[name] == selected {
					values = append(values, item)
				}
			}
		} else if isObject {
			values = []any{object[segments[0]]}
		}
	}
	return !slices.ContainsFunc(values, func(item any) bool { return !hasValueType(item, segments[1:], valueType) })
}

// valueExample shows a native value for valuePath, for example [123] for a list of integers, or jsondecode(...) for an
// object written as a JSON string. Lists of strings have no example.
func valueExample(valuePath string, valueTypes map[string]string) string {
	switch valueTypes[valuePath] {
	case "integer":
		return "123"
	case "boolean":
		return "true"
	case "object":
		return "jsondecode(...)"
	case "list":
		if item := valueExample(valuePath+"[]", valueTypes); item != "" {
			return "[" + item + "]"
		}
	}
	return ""
}

func JSONValueToAPI(ctx context.Context, attributePath tfpath.Path, source attr.Value) (any, diag.Diagnostics) {
	if source.IsNull() || source.IsUnknown() {
		return nil, nil
	}
	value, diags := attributeToInterface(ctx, attributePath, source, false, true)
	if diags.HasError() || value == nil {
		return value, diags
	}
	if encoded, ok := value.(string); ok {
		if err := json.Unmarshal([]byte(encoded), &value); err != nil {
			diags.AddAttributeError(attributePath, "Invalid JSON value", "Expected valid JSON or a native HCL value. Remove jsonencode(...) when migrating. The JSON string could not be decoded.")
			return nil, diags
		}
	}
	return value, diags
}

func JSONToStringMapSlice(attributePath tfpath.Path, value any) ([]map[string]any, diag.Diagnostics) {
	if value == nil {
		return nil, nil
	}
	items, ok := value.([]any)
	if !ok {
		return nil, diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid JSON value", "Expected an array of objects.")}
	}
	result := make([]map[string]any, len(items))
	for index, item := range items {
		object, ok := item.(map[string]any)
		if !ok {
			return nil, diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath.AtListIndex(index), "Invalid JSON value", "Expected an object.")}
		}
		result[index] = object
	}
	return result, nil
}

func APIToDynamicJSON(ctx context.Context, attributePath tfpath.Path, source any, prior types.Dynamic, mapPaths, writeOnly []string, normalization string) (types.Dynamic, diag.Diagnostics) {
	// SDK array properties use []map[string]any; decode once into the same JSON tree as configuration.
	encodedSource, err := json.Marshal(source)
	if err != nil {
		return types.DynamicNull(), diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid API value", "Could not encode the API response as JSON.")}
	}
	if err := json.Unmarshal(encodedSource, &source); err != nil {
		return types.DynamicNull(), diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid API value", "Could not decode the API response as JSON.")}
	}
	if prior.IsNull() || prior.IsUnknown() || prior.IsUnderlyingValueNull() || prior.IsUnderlyingValueUnknown() {
		return ToDynamic(ctx, attributePath, source, nil)
	}
	configured, diags := JSONValueToAPI(ctx, attributePath, prior)
	if diags.HasError() {
		return types.DynamicNull(), diags
	}
	// Email and folder targets store a null delivery_policy as the default batch_interval. Other target types store {}.
	if normalization == "event_target_delivery_policy" && configured == nil && reflect.DeepEqual(source, map[string]any{"batch_interval": batchIntervalMin}) {
		source = map[string]any{}
	}
	if normalization != "" {
		source = preserveJSONNormalization(source, configured, normalizedJSON(configured, normalization), normalization == "mount_mappings")
	}
	value, _ := configuredDynamicStateValue(source, source != nil, configured, writeOnly, "", mapPaths...)
	if reflect.DeepEqual(configured, value) {
		return prior, nil
	}
	if _, encoded := prior.UnderlyingValue().(types.String); encoded {
		result, err := json.Marshal(value)
		if err != nil {
			return types.DynamicNull(), diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid API value", "Could not encode the API value as a JSON string.")}
		}
		return types.DynamicValue(types.StringValue(string(result))), nil
	}
	result, diags := dynamicJSONAttribute(ctx, attributePath, value, prior.UnderlyingValue(), mapPaths, "")
	if diags.HasError() {
		return types.DynamicNull(), diags
	}
	return types.DynamicValue(result), nil
}

// Keep equivalent configured spellings, but still record changes to other fields.
func preserveJSONNormalization(source, configured, normalized any, normalizedKeys bool) any {
	if reflect.DeepEqual(source, normalized) {
		return configured
	}
	switch prior := configured.(type) {
	case map[string]any:
		values, ok := source.(map[string]any)
		canonical, canonicalObject := normalized.(map[string]any)
		if !ok || !canonicalObject {
			return source
		}
		for key, value := range prior {
			canonicalKey := key
			if normalizedKeys {
				canonicalKey = strings.TrimSpace(key)
			}
			current, exists := values[canonicalKey]
			if exists || canonical[canonicalKey] == nil {
				values[key] = preserveJSONNormalization(current, value, canonical[canonicalKey], false)
				if key != canonicalKey {
					delete(values, canonicalKey)
				}
			}
		}
	case []any:
		values, ok := source.([]any)
		canonical, canonicalList := normalized.([]any)
		if ok && canonicalList {
			for index := range values {
				if index < len(prior) && index < len(canonical) {
					values[index] = preserveJSONNormalization(values[index], prior[index], canonical[index], false)
				}
			}
		}
	}
	return source
}

func dynamicJSONAttribute(ctx context.Context, attributePath tfpath.Path, source any, prior attr.Value, mapPaths []string, currentPath string) (attr.Value, diag.Diagnostics) {
	if prior == nil {
		value, diags := ToDynamic(ctx, attributePath, source, nil)
		if value.IsNull() {
			return types.DynamicNull(), diags
		}
		return value.UnderlyingValue(), diags
	}
	if dynamic, ok := prior.(types.Dynamic); ok {
		value, diags := dynamicJSONAttribute(ctx, attributePath, source, dynamic.UnderlyingValue(), mapPaths, currentPath)
		return types.DynamicValue(value), diags
	}
	if source == nil {
		value, err := prior.Type(ctx).ValueFromTerraform(ctx, tftypes.NewValue(prior.Type(ctx).TerraformType(ctx), nil))
		if err != nil {
			return nil, diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid API value", err.Error())}
		}
		return value, nil
	}
	// An API value of another kind than prior, such as a list where a different Folder Behavior had an object, has no
	// configured representation to keep. It is read as the API sends it, so the change shows as drift.
	if object, ok := prior.(types.Object); ok {
		values, ok := source.(map[string]any)
		if !ok {
			return dynamicJSONAttribute(ctx, attributePath, source, nil, mapPaths, currentPath)
		}
		if !apiKeyedPath(mapPaths, currentPath) {
			for name := range object.AttributeTypes(ctx) {
				if _, exists := values[name]; !exists {
					values[name] = nil
				}
			}
		}
		attributes := make(map[string]attr.Value, len(values))
		attributeTypes := make(map[string]attr.Type, len(values))
		var diags diag.Diagnostics
		for name, value := range values {
			entryPath := name
			if currentPath != "" {
				entryPath = currentPath + "." + name
			}
			entry, entryDiags := dynamicJSONAttribute(ctx, attributePath.AtName(name), value, object.Attributes()[name], mapPaths, entryPath)
			diags.Append(entryDiags...)
			if !entryDiags.HasError() {
				attributes[name], attributeTypes[name] = entry, entry.Type(ctx)
			}
		}
		if diags.HasError() {
			return nil, diags
		}
		return types.ObjectValue(attributeTypes, attributes)
	}
	if tuple, ok := prior.(types.Tuple); ok {
		values, ok := source.([]any)
		if !ok {
			return dynamicJSONAttribute(ctx, attributePath, source, nil, mapPaths, currentPath)
		}
		elements := make([]attr.Value, len(values))
		elementTypes := make([]attr.Type, len(values))
		var diags diag.Diagnostics
		for index, value := range values {
			var old attr.Value
			if index < len(tuple.Elements()) {
				old = tuple.Elements()[index]
			}
			entry, entryDiags := dynamicJSONAttribute(ctx, attributePath.AtListIndex(index), value, old, mapPaths, currentPath)
			diags.Append(entryDiags...)
			if !entryDiags.HasError() {
				elements[index], elementTypes[index] = entry, entry.Type(ctx)
			}
		}
		if diags.HasError() {
			return nil, diags
		}
		return types.TupleValue(elementTypes, elements)
	}
	if collection, ok := prior.(types.Map); ok {
		values, ok := source.(map[string]any)
		if !ok {
			return dynamicJSONAttribute(ctx, attributePath, source, nil, mapPaths, currentPath)
		}
		elements := make(map[string]attr.Value, len(values))
		elementTypes := make(map[string]attr.Type, len(values))
		var diags diag.Diagnostics
		for name, value := range values {
			entryPath := name
			if currentPath != "" {
				entryPath = currentPath + "." + name
			}
			entry, entryDiags := dynamicJSONAttribute(ctx, attributePath.AtMapKey(name), value, collection.Elements()[name], mapPaths, entryPath)
			diags.Append(entryDiags...)
			if !entryDiags.HasError() {
				elements[name], elementTypes[name] = entry, entry.Type(ctx)
			}
		}
		if diags.HasError() {
			return nil, diags
		}
		if uniformElementType(ctx, collection.ElementType(ctx), slices.Collect(maps.Values(elements))) {
			return types.MapValue(collection.ElementType(ctx), elements)
		}
		return types.ObjectValue(elementTypes, elements)
	}
	if list, ok := prior.(types.List); ok {
		return dynamicJSONCollection(ctx, attributePath, source, list.Elements(), list.ElementType(ctx), false, mapPaths, currentPath)
	}
	if set, ok := prior.(types.Set); ok {
		return dynamicJSONCollection(ctx, attributePath, source, set.Elements(), set.ElementType(ctx), true, mapPaths, currentPath)
	}
	switch prior.(type) {
	case types.Number, types.Int64, types.Float64, types.String, types.Bool:
		value, diags := ToDynamic(ctx, attributePath, source, prior)
		return value.UnderlyingValue(), diags
	}
	return schemaValue(ctx, attributePath, source, prior.Type(ctx))
}

// Typed lists and sets keep their element type while every element still matches it. Otherwise the value becomes a tuple so a changed API type shows as a diff.
func dynamicJSONCollection(ctx context.Context, attributePath tfpath.Path, source any, priorElements []attr.Value, elementType attr.Type, isSet bool, mapPaths []string, currentPath string) (attr.Value, diag.Diagnostics) {
	values, ok := source.([]any)
	if !ok {
		return dynamicJSONAttribute(ctx, attributePath, source, nil, mapPaths, currentPath)
	}
	elements := make([]attr.Value, len(values))
	elementTypes := make([]attr.Type, len(values))
	var diags diag.Diagnostics
	for index, value := range values {
		var old attr.Value
		if index < len(priorElements) {
			old = priorElements[index]
		}
		entry, entryDiags := dynamicJSONAttribute(ctx, attributePath.AtListIndex(index), value, old, mapPaths, currentPath)
		diags.Append(entryDiags...)
		if !entryDiags.HasError() {
			elements[index], elementTypes[index] = entry, entry.Type(ctx)
		}
	}
	if diags.HasError() {
		return nil, diags
	}
	if !uniformElementType(ctx, elementType, elements) {
		return types.TupleValue(elementTypes, elements)
	}
	if isSet {
		return types.SetValue(elementType, elements)
	}
	return types.ListValue(elementType, elements)
}

func uniformElementType(ctx context.Context, elementType attr.Type, elements []attr.Value) bool {
	for _, element := range elements {
		if !element.Type(ctx).Equal(elementType) {
			return false
		}
	}
	return true
}

// CurrentFormFields returns the fields listed in form_layout, in that order. Authenticated Form Field Set
// responses also contain historical definitions that existing registration answers use. fields is the form_fields
// value of the response body, as ResponseProperty returns it.
func CurrentFormFields(formFieldSet files_sdk.FormFieldSet, fields any) any {
	all, _ := fields.([]any)
	fieldsByID := make(map[float64]map[string]any, len(all))
	for _, field := range all {
		if record, ok := field.(map[string]any); ok {
			if id, ok := record["id"].(float64); ok {
				fieldsByID[id] = record
			}
		}
	}
	current := make([]any, 0, len(formFieldSet.FormLayout))
	for _, id := range formFieldSet.FormLayout {
		if field, ok := fieldsByID[float64(id)]; ok {
			current = append(current, field)
		}
	}
	return current
}

// AlignMembersByEmail puts Share Group members in the order and email spelling of prior, the configured or prior
// state value. Rails matches members by email ignoring letter case, keeps the stored spelling and returns them in its
// own order. Unmatched members follow in API order. source is the members value of the response body.
func AlignMembersByEmail(ctx context.Context, source any, prior attr.Value) any {
	configured, _ := JSONValueToAPI(ctx, tfpath.Empty(), prior)
	members, isList := source.([]any)
	entries, configuredList := configured.([]any)
	if !isList || !configuredList {
		return source
	}
	remaining := slices.Clone(members)
	aligned := make([]any, 0, len(members))
	for _, entry := range entries {
		object, _ := entry.(map[string]any)
		email, _ := object["email"].(string)
		index := slices.IndexFunc(remaining, func(member any) bool {
			memberObject, _ := member.(map[string]any)
			memberEmail, _ := memberObject["email"].(string)
			return email != "" && strings.EqualFold(memberEmail, email)
		})
		if index < 0 {
			continue
		}
		member := remaining[index].(map[string]any)
		member["email"] = email
		aligned = append(aligned, member)
		remaining = slices.Delete(remaining, index, index+1)
	}
	return append(aligned, remaining...)
}

// apiKeyedPath reports whether the API decides the keys at currentPath, so keys it adds stay visible as drift.
// "a.b" is a map at that path. "a.b.*" is a free-form value covering every path from a.b down, and "*" covers the whole value.
func apiKeyedPath(mapPaths []string, currentPath string) bool {
	return slices.ContainsFunc(mapPaths, func(mapPath string) bool {
		if root, freeForm := strings.CutSuffix(mapPath, "*"); freeForm {
			root = strings.TrimSuffix(root, ".")
			return root == "" || currentPath == root || strings.HasPrefix(currentPath, root+".")
		}
		return mapPath == currentPath
	})
}

// Compare the documented Rails normalizations without rewriting requests or hiding unrelated drift.
// Each case mirrors a Rails method in files-rails. Change both sides together.
//
//	"mount_mappings": DesktopConfigurationProfile#normalize_mount_mappings and Path.normalize
//	"expected_remote_servers": IntegrationCentricProfile#normalize_expected_remote_server
//	"watermark": FolderBehavior::Watermark.normalize_integer_keys
//	"empty_object": json_property columns that store {} when the request sends null
//	"event_target_config": EventTargets::DeliveryAdapters::GooglePubSub#normalized_config and #sanitized_config
//	"event_target_delivery_policy": EventTarget#normalize_json_attribute and #normalize_batch_delivery_policy
//	"form_fields": FormField#clear_unrelated_fields and the boolean cast of required
//
// Managed reads of related records use the response body, so an empty string and null stay distinct.
func normalizedJSON(source any, normalization string) any {
	encoded, _ := json.Marshal(source)
	var value any
	if json.Unmarshal(encoded, &value) != nil {
		return source
	}
	if value == nil && (normalization == "empty_object" || normalization == "watermark" || normalization == "event_target_delivery_policy") {
		return map[string]any{}
	}
	switch normalization {
	case "event_target_config":
		if object, ok := value.(map[string]any); ok {
			if text, ok := object["google_credentials"].(string); ok {
				var credentials map[string]any
				if json.Unmarshal([]byte(text), &credentials) == nil && credentials != nil {
					delete(credentials, "private_key")
					object["google_credentials"] = credentials
				}
			}
		}
	case "event_target_delivery_policy":
		if object, ok := value.(map[string]any); ok {
			if interval, exists := object["batch_interval"]; exists {
				// Email and Folder targets store EventTarget::BATCH_INTERVAL_MIN for a blank value and an integer for digits.
				// Other target types keep the value, so it still compares equal to the configuration.
				if text, ok := interval.(string); ok && batchIntervalDigits.MatchString(text) {
					number, _ := strconv.ParseInt(text, 10, 64)
					object["batch_interval"] = float64(number)
				} else if jsonBlank(interval) {
					object["batch_interval"] = batchIntervalMin
				}
			}
		}
	case "form_fields":
		if entries, ok := value.([]any); ok {
			for _, entry := range entries {
				object, ok := entry.(map[string]any)
				if !ok {
					continue
				}
				if fieldType, _ := object["field_type"].(string); fieldType == "text" || fieldType == "text_area" {
					delete(object, "options_for_select")
					delete(object, "default_option")
				}
				// FormField#required is a Rails boolean column: it stores these strings as false and other text as true.
				if text, ok := object["required"].(string); ok && text != "" {
					object["required"] = !slices.Contains(railsFalseStrings, text)
				}
			}
		}
	case "mount_mappings":
		if object, ok := value.(map[string]any); ok {
			result := make(map[string]any, len(object))
			for key, value := range object {
				if text, ok := value.(string); ok {
					parts := strings.FieldsFunc(strings.TrimSpace(text), func(char rune) bool { return char == '/' || char == '\\' })
					cleaned := parts[:0]
					for _, part := range parts {
						part = strings.ReplaceAll(part, "\x00", "")
						if part != "" && part != "." && part != ".." {
							cleaned = append(cleaned, part)
						}
					}
					value = strings.Join(cleaned, "/")
				}
				result[strings.TrimSpace(key)] = value
			}
			return result
		}
	case "expected_remote_servers":
		if entries, ok := value.([]any); ok {
			for _, entry := range entries {
				if object, ok := entry.(map[string]any); ok {
					switch name := object["name"].(type) {
					case float64:
						object["name"] = strconv.FormatFloat(name, 'f', -1, 64)
					case bool:
						object["name"] = strconv.FormatBool(name)
					}
					for _, key := range []string{"name", "server_type"} {
						if text, ok := object[key].(string); ok {
							object[key] = strings.TrimSpace(text)
						}
					}
					if object["name"] == "" {
						delete(object, "name")
					}
				}
			}
		}
	case "watermark":
		if object, ok := value.(map[string]any); ok {
			for _, key := range []string{"transparency", "max_height_or_width"} {
				if text, ok := object[key].(string); ok {
					if number, err := strconv.ParseInt(strings.TrimSpace(text), 10, 64); err == nil {
						object[key] = float64(number)
					}
				}
			}
		}
	}
	return value
}

var batchIntervalDigits = regexp.MustCompile(`^[0-9]+$`)

// railsFalseStrings are the strings that ActiveModel::Type::Boolean casts to false.
var railsFalseStrings = []string{"0", "f", "F", "false", "FALSE", "off", "OFF"}

// batchIntervalMin mirrors EventTarget::BATCH_INTERVAL_MIN.
const batchIntervalMin = float64(600)

// jsonBlank mirrors Ruby's blank? for decoded JSON values.
func jsonBlank(value any) bool {
	switch value := value.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(value) == ""
	case bool:
		return !value
	case []any:
		return len(value) == 0
	case map[string]any:
		return len(value) == 0
	}
	return false
}
