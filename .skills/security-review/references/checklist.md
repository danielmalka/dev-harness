# Trust boundary checklist

Use this after step 1 of the procedure, once the boundaries and exposures in scope are known. Each row names what to look for and the evidence needed to classify it. Demonstrated vulnerabilities need a reachable unmitigated path; confirmed exposures need exposure evidence, without requiring an application exploit chain. Keep incomplete evidence qualified as a hypothesis. A row that does not apply to the change is skipped, not answered.

## Identity and access

| Area | Look for | Trace before reporting |
| --- | --- | --- |
| Authentication | New or moved entry points that skip the session check; token validation that ignores expiry, audience, or signature; credentials compared without a constant-time function | The exact route or handler that reaches sensitive work without a verified identity |
| Authorization | A route guard with no object-level check; a role read from a client-supplied value; a permission computed after the data was already loaded and returned; a check that does not fail closed on a missing or partial identity; cache or index keys not fully scoped to tenant, user, and object | Which caller can name which object identifier, and where ownership is enforced |
| Session | Tokens in URLs or logs; missing rotation after privilege change; cookies without the flags the project's transport requires | How a stolen or fixed session value reaches an authenticated action |
| Privilege boundaries | Internal endpoints reachable from outside; administrative flags derived from user input; impersonation without an audit record | The shortest path from an unprivileged position to the privileged operation |

## Input and interpretation

| Area | Look for | Trace before reporting |
| --- | --- | --- |
| Query construction | String concatenation or interpolation into SQL or another query language; a dynamic identifier such as a table, column, or sort field | Whether every path to the sink uses parameters, and what the non-parameterized branch accepts |
| Command execution | Shell invocation with interpolated arguments; a shell used where a direct process call would do; arguments assembled from configuration a user can set | Which characters survive to the shell and what they can start |
| Path handling | Filenames from requests or archives joined to a base directory; traversal sequences; symbolic links followed on extraction | Whether the resolved path is confirmed to stay inside the intended root after normalization |
| Deserialization | Object graphs rebuilt from untrusted bytes; formats that can name a type; unsafe YAML or pickle equivalents | What types the payload can instantiate and what their construction executes |
| Template rendering | User data passed as template source rather than template data; escaping disabled for convenience | Whether the escape is on by default on the path that renders the untrusted value |
| Model and tool boundary | Untrusted text that reaches a model, or a tool with authority, is a trust boundary; tool policy enforced inside the prompt instead of outside the model; retrieved text or tool output treated as instructions instead of data | Where tool policy is enforced and whether the model can change it |
| Output encoding | Values written into HTML, JSON embedded in HTML, attributes, or URLs without the encoding that context requires | The exact rendering context and the encoder applied there |

## Data and secrets

| Area | Look for | Trace before reporting |
| --- | --- | --- |
| Secrets in the tree | Keys, tokens, passwords, connection strings in source, fixtures or configuration; secrets in error text. Scanning version-control history requires a shell; when unavailable, record history as not inspected | The file and line. Report the location and shape; never the value. A leaked secret is rotated or revoked outside git before it is dismissed or baselined; a baseline holds only exact, reviewed fingerprints, never a broad exclusion by path, rule, or pattern |
| Sensitive data in transit | Confidential fields serialized into responses, logs, metrics, traces, or analytics events | Which consumers receive the field and whether they are entitled to it |
| Storage | Confidential values stored without the protection the project requires; backups or exports that widen access | Where the value lands and who can read that location |
| Cryptography | Home-made algorithms; fixed or predictable initialization vectors and salts; a hash where a password-hashing function belongs; a random source not meant for security | Which property the code depends on and where the weak primitive breaks it |
| Retention and deletion | Data kept past its stated purpose; delete operations that leave copies behind | Whether the documented deletion actually removes every copy the change creates |

## Outbound and dependencies

| Area | Look for | Trace before reporting |
| --- | --- | --- |
| Outbound requests | URLs built from user input; redirects followed without restriction; requests to hosts resolved at call time | Whether an attacker can aim the request at an internal address or a metadata service |
| Redirects | Destination taken from a request parameter with no allowlist | Which values reach the redirect and what the user's browser then trusts |
| Third-party responses | Data from an external service treated as trusted; webhooks accepted without signature verification | Where the external value becomes a decision or a sink argument |
| Dependencies | Newly added packages; versions pinned to something unmaintained; a transitive change that widens capability; lockfile drift against the manifest; new registries or git dependencies; package lifecycle hooks (install, postinstall) | What the new dependency is allowed to do in the process |
| Configuration | Debug or verbose modes enabled by default; permissive origins; disabled certificate verification; wide file permissions; a scanner config, ignore list, baseline, or threshold that the diff narrowed or made non-gating (a clean result is then not a pass) | Whether the risky setting is what the code ships with, not just what a deployment could choose |
| Hostile or unexplained content | Executable bits, symlinks, submodules, Unicode bidi characters or homoglyphs, encoded or binary blobs with no stated purpose. CI workflow changes: `pull_request_target`, widened `permissions`, unpinned actions, secrets that reach untrusted code | What the file does when it runs or is checked out, and who can trigger it |
| Weakened controls | The diff disables or weakens a security control (tests, logging, TLS or signature checks, auth, limits, scanners); tests or mocks that hide real side effects; delayed or conditional activation by time, host, CI, or region | Whether the control still runs on the path an attacker uses, and what the mock replaced |

## Resilience and observability

| Area | Look for | Trace before reporting |
| --- | --- | --- |
| Error handling | Stack traces or internal identifiers returned to callers; errors swallowed so a failed security check looks like success | Whether a failure on the security path produces an allow |
| Rate and size limits | Unbounded uploads, loops, or expansions driven by attacker input; no limit on authentication attempts | The input that makes cost grow and by how much |
| Concurrency | A check and its use separated by a window another request can exploit | The interleaving that produces the bypass |
| Audit | Security-relevant actions with no record; logs an attacker can forge by injecting separators into a field | Which action becomes untraceable and who can make it so |

## Proportionality

The checklist is a prompt, not a form to fill in. A change that adds one field to an internal response touches one or two rows. A change that adds an upload endpoint touches many. Report only what you traced, and say which rows you skipped and why.
