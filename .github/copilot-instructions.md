# Copilot Instructions for Evilginx (chromeworm)

## Project Overview
- Evilginx is a Go-based man-in-the-middle attack framework for phishing credentials and session cookies, bypassing 2FA.
- The core architecture includes a custom HTTP(S) proxy, DNS server, and a modular phishlet system for targeting specific sites.
- The project is structured for extensibility: new phishlets (YAML configs) define site-specific proxying and rewriting logic.

## Key Components
- `main.go`: Entry point, CLI, and orchestrator.
- `core/`: Main logic for proxying, phishlet parsing, session handling, and utilities.
  - `http_proxy.go`: HTTP(S) proxy logic, request/response rewriting, session management.
  - `phishlet.go`: Phishlet config parsing, Go struct mapping, and phishlet-specific logic.
  - `config.go`, `session.go`, `utils.go`: Configuration, session, and helper utilities.
- `phishlets/`: YAML files describing how to proxy and rewrite specific target sites.
- `database/`: Session and credential storage.
- `log/`: Logging utilities.
- `media/`: Images for documentation/UI.

## Developer Workflows
- **Build:** Use `build.bat` (Windows) or `make` (if Makefile is present) to build the binary.
- **Run:** Use `build_run.bat` or run the built binary directly. Requires admin/root for DNS/HTTP(S) binding.
- **Phishlet Development:** Add or modify YAML files in `phishlets/`. Update Go structs in `phishlet.go` if schema changes.
- **Proxy/Rewrite Logic:** Extend `http_proxy.go` for new request/response handling features.

## Project-Specific Patterns
- **Phishlet Parsing:** All YAML config fields must be explicitly mapped in `ConfigPhishlet` and related structs in `phishlet.go`.
- **Session Mapping:** Session and token tracking is handled in-memory and/or via the `database/` package.
- **URL Rewriting:** Use the `rewrite_urls` section in phishlets and corresponding Go logic to obfuscate paths/params.
- **Error Handling:** Most errors are surfaced to the CLI; phishlet parsing errors are explicit and descriptive.

## External Dependencies
- Uses several Go libraries: `goproxy`, `viper`, `mux`, `lego`, etc. (see `go.mod`).
- Phishlet YAMLs are not bundled—users must supply or author them.

## Example: Adding a New Phishlet Field
1. Add the field to the YAML in `phishlets/example.yaml`.
2. Update `ConfigPhishlet` and related structs in `phishlet.go`.
3. Add parsing/logic in `LoadFromFile` and proxy handling in `http_proxy.go` if needed.

## References
- [Official Evilginx Documentation](https://help.evilginx.com)
- [Evilginx Mastery Course](https://academy.breakdev.org/evilginx-mastery)
- [Project README](../README.md)

---
If any section is unclear or missing, please provide feedback for improvement.
