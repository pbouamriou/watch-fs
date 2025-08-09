# NEXT-FEATURES

This document tracks proposed features and their progress. Priorities are ordered by value first, then effort. Status codes: [ ] Planned, [WIP] In Progress, [x] Done.

Priority legend:

- Value: High / Medium / Low
- Effort: S (small) / M (medium) / L (large)

## Top Priorities (Highest Value)

1. Event Rules & Actions (Alerts/Hooks)

- Value: High | Effort: M
- Description: Define rules (match by path, op, type, regex) and trigger actions: shell command, HTTP webhook, desktop notification.
- Acceptance:
  - Rules file (YAML) and CLI flags to enable
  - Supports includes/excludes, op filters, and regex on path
  - Actions: command, webhook; minimal retry & backoff for webhooks
  - TUI: minimal indicator when a rule triggers (optional initial)
- Tasks:
  - Rule schema + parser
  - Matcher + evaluator
  - Action runners (command/webhook)
  - Unit tests + examples

2. Include/Exclude Patterns (+ .gitignore / .watchfsignore)

- Value: High | Effort: M
- Description: Noise reduction through globs/regex and repo-aware ignores.
- Acceptance:
  - `--include`, `--exclude` flags (multi)
  - Respect `.gitignore` and optional `.watchfsignore` per root
  - Hot-reload ignores on file change (nice-to-have)
- Tasks: pattern engine, precedence rules, fsnotify ignore handling, tests

3. Real-time Search/Filter in TUI

- Value: High | Effort: S-M
- Description: Quick filter (e.g. `/` to search), supports substring and regex, highlights matches.
- Acceptance:
  - Open search box with `/`, close with ESC
  - Live filter of events list without data loss
  - Clear visual highlight of matched substrings
- Tasks: input mode, filter pipeline, highlighting, tests

4. Pause/Resume per Root (and Global)

- Value: High | Effort: S
- Description: Temporarily pause watching to reduce noise or focus on subsets.
- Acceptance:
  - Keybinding to toggle pause for selected root; global pause
  - Clear TUI indicator (icon/color) and excluded-from-count behavior
- Tasks: watcher gating, UI state, help text, tests

5. Config File & Profiles

- Value: High | Effort: M
- Description: Centralized YAML config for paths, filters, sort, keybindings, theme. Named profiles for different projects.
- Acceptance:
  - `-config` flag; load default from `$XDG_CONFIG_HOME/watch-fs/config.yml`
  - Select profile with `--profile` or via TUI
  - Safe merging of CLI flags over config
- Tasks: schema, loader, precedence, docs, tests

6. Persistent History with Rolling Store

- Value: High | Effort: M-L
- Description: Continuously persist events to SQLite with size/time-based rotation; reload last session.
- Acceptance:
  - Background writer with backpressure
  - Rotation by file size and/or age
  - TUI option to reopen previous session
- Tasks: writer, rotation policy, recovery, tests

## Next (High/Medium Value)

- Stats/Timeline View

  - Value: Medium-High | Effort: M
  - Sparkline or mini chart of events/minute, per operation and per root.

- Server Mode + HTTP API (SSE/WebSocket)

  - Value: High | Effort: L
  - Headless mode streaming events; minimal web UI to browse/filter.

- Export/Import Enhancements (CSV, filtered export)

  - Value: Medium | Effort: S-M
  - CSV export, selective export by current filter, scheduled export.

- Theme & Keybinding Customization

  - Value: Medium | Effort: S-M
  - Custom themes (colors), override shortcuts via config.

- Quick Diff/Tail for Text Files
  - Value: Medium | Effort: M
  - From an event, open a small popup to show last lines or external diff.

## Later (Nice-to-have)

- Desktop notifications integration (macOS/Linux/Windows)
- Saved queries (named filters) and bookmarks/tags on events
- Multi-tab workspaces (per root or per filter)
- i18n of TUI strings

---

## Tracking Checklist

- [ ] Event Rules & Actions (Alerts/Hooks)
- [ ] Include/Exclude Patterns (+ .gitignore)
- [ ] Real-time Search/Filter in TUI
- [ ] Pause/Resume per Root (and Global)
- [ ] Config File & Profiles
- [ ] Persistent History with Rolling Store
- [ ] Stats/Timeline View
- [ ] Server Mode + HTTP API (SSE/WebSocket)
- [ ] Export/Import Enhancements (CSV, filtered export)
- [ ] Theme & Keybinding Customization
- [ ] Quick Diff/Tail for Text Files

## Notes

- Keep code style consistent with Go guidelines and existing architecture (composition, focus management, thread safety).
- Prefer interfaces for new subsystems (rules engine, action runners) for testability.
- Always add tests and update help/usage in README when user-facing changes are introduced.
