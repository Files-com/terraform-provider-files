resource "files_ai_assistant_personality" "example_ai_assistant_personality" {
  apply_to_all_workspaces = false
  system_prompt           = "Respond as a concise operations assistant."
  use_by_default          = false
  workspace_id            = 0
}

