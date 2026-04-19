#!/usr/bin/env bash
# Runs every Go fuzz target in gopptx for FUZZ_TIME each.
# Crash inputs that the fuzzer finds are written to
#   $FUZZ_OUT/<binary>/<FuzzName>/
# for easy inspection after the container exits.
#
# Usage (inside Docker):
#   FUZZ_TIME=30s FUZZ_OUT=/workspace/fuzz-out ./scripts/fuzz.sh
#
# Usage (locally, from repo root):
#   FUZZ_TIME=10s ./scripts/fuzz.sh

set -euo pipefail

FUZZ_TIME="${FUZZ_TIME:-30s}"
FUZZ_OUT="${FUZZ_OUT:-$(pwd)/fuzz-out}"

# Map of compiled binary path → space-separated fuzz function names.
# The binaries are pre-built into /usr/local/bin by Dockerfile.fuzz so that
# the fuzz loop needs no write access to the Go module/build cache.
declare -A TARGETS=(
  ["/usr/local/bin/fuzz-pptxxml"]="FuzzEscape FuzzContentTypes FuzzSectionListXML FuzzPresentation FuzzEmbeddedFontsXML FuzzRichTextRun"
  ["/usr/local/bin/fuzz-markdown"]="FuzzSlidesFromMarkdown FuzzParseInlineTextRuns"
  ["/usr/local/bin/fuzz-netsec"]="FuzzValidateURLForHTTPFetch FuzzIsBlockedAddr"
  ["/usr/local/bin/fuzz-urlfetch"]="FuzzWebParserParse FuzzCheckAddrBlocked"
  ["/usr/local/bin/fuzz-tplx"]="FuzzMergeAdjacentRuns FuzzInterpolateText"
)

PASS=0
FAIL=0
CRASH=0

run_fuzz() {
  local bin="$1"
  local fn="$2"
  local label
  label="$(basename "$bin"):${fn}"
  local crash_dir="${FUZZ_OUT}/${label//:/_}"

  mkdir -p "$crash_dir"

  echo "━━━ ${label}  [${FUZZ_TIME}] ━━━"

  if [[ ! -x "$bin" ]]; then
    echo "  SKIP (binary not found: ${bin})"
    return
  fi

  # Run the pre-compiled test binary directly.
  # -test.fuzz       — which fuzz function to run
  # -test.fuzztime   — wall-clock budget
  # -test.fuzzminimizetime=0s — skip minimisation to keep the run bounded
  # -test.fuzzcachedir — write the generated corpus here so crashes are visible
  if "$bin" \
      -test.fuzz="^${fn}$" \
      -test.fuzztime="${FUZZ_TIME}" \
      -test.fuzzminimizetime=0s \
      -test.fuzzcachedir="${crash_dir}/corpus" 2>&1; then
    echo "  PASS"
    PASS=$((PASS + 1))
  else
    echo "  FAIL (possible crash or error — check ${crash_dir})"
    FAIL=$((FAIL + 1))
    # Count crash files (corpus entries that triggered a failure).
    local n
    n=$(find "$crash_dir" -type f 2>/dev/null | wc -l)
    if [[ "$n" -gt 0 ]]; then
      echo "  ${n} file(s) saved to ${crash_dir}"
      CRASH=$((CRASH + 1))
    fi
  fi
}

echo "gopptx fuzz run  |  budget=${FUZZ_TIME}  |  output=${FUZZ_OUT}"
echo

for bin in "${!TARGETS[@]}"; do
  for fn in ${TARGETS[$bin]}; do
    run_fuzz "$bin" "$fn"
  done
done

echo
echo "━━━ Summary ━━━"
echo "  Passed : ${PASS}"
echo "  Failed : ${FAIL}"
echo "  Crashes: ${CRASH}"

if [[ "$CRASH" -gt 0 ]]; then
  echo "  Crash inputs written to: ${FUZZ_OUT}"
  exit 1
fi
