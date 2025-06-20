#!/usr/bin/env bash

PRODUCTS_FILE="products.json"
OS_FILE="os.json"
LOG_FILE="domain_mismatches.log"

COMPARE=true
DEBUG=false

while [[ $# -gt 0 ]]; do
  case $1 in
    --no-compare)
      COMPARE=false
      shift
      ;;
    --debug)
      DEBUG=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--no-compare] [--debug]"
      exit 1
      ;;
  esac
done

if ! command -v jq &> /dev/null; then
    echo "jq is required but not installed. Exiting."
    exit 1
fi

> "$LOG_FILE"

# Map each environment label (e.g., prod, staging) to its respective bouncer domain
LABELS=("prod" "stage")
DOMAINS=("download.mozilla.org" "stage.bouncer.nonprod.webservices.mozgcp.net")

PRODUCTS=$(jq -c '.[]' "$PRODUCTS_FILE")
OS_LIST=$(jq -r '.[].name' "$OS_FILE")

echo "🔍 Starting consistency check across domains..."
$COMPARE || echo "⚠️  Domain comparison is DISABLED (--no-compare)"
$DEBUG && echo "* Debug mode is ENABLED"

for product in $PRODUCTS; do
    alias=$(echo "$product" | jq -r '.alias')
    lang=$(echo "$product" | jq -r '.language')

    for os in $OS_LIST; do
        echo "🔎 Checking: product=$alias, os=$os, lang=$lang"

        PATHS=()
        LOCATIONS=()

        for ((i=0; i<${#DOMAINS[@]}; i++)); do
            label=${LABELS[$i]}
            domain=${DOMAINS[$i]}
            url="https://${domain}/?product=${alias}&os=${os}&lang=${lang}"

            $DEBUG && echo "➡️  CURL [$label]: curl -s -v \"$url\""

            # Capture verbose curl output (headers + GET line)
            response=$(curl -s -v "$url" 2>&1)

            # Extract GET line path
            get_path=$(echo "$response" | grep -Eo 'GET .* HTTP' | sed -E 's/GET (.*) HTTP.*/\1/')
            PATHS[$i]="$get_path"

            # Extract location header (case-insensitive match)
            location=$(echo "$response" | grep -i '^< location:' | sed -E 's/^< location: *//I')
            LOCATIONS[$i]="$location"

            $DEBUG && echo "   ↪️  Path:     $get_path"
            $DEBUG && echo "   ↪️  Location: $location"
        done

        # Compare logic
        if $COMPARE; then
            ref_path="${PATHS[0]}"
            ref_loc="${LOCATIONS[0]}"

            for ((i=1; i<${#LABELS[@]}; i++)); do
                label=${LABELS[$i]}
                path="${PATHS[$i]}"
                loc="${LOCATIONS[$i]}"

                if [[ "$path" != "$ref_path" ]]; then
                    echo -e "\033[1;91mX\033[0m PATH mismatch for $alias (lang=$lang, os=$os): $label path differs from prod" | tee -a "$LOG_FILE"
                    echo "    $label path: $path" | tee -a "$LOG_FILE"
                    echo "    prod path : $ref_path" | tee -a "$LOG_FILE"
                    echo "---" >> "$LOG_FILE"
                fi

                if [[ "$loc" != "$ref_loc" ]]; then
                    echo -e "\033[1;91mX\033[0m LOCATION mismatch for $alias (lang=$lang, os=$os): $label location differs from prod" | tee -a "$LOG_FILE"
                    echo "    $label location: $loc" | tee -a "$LOG_FILE"
                    echo "    prod location : $ref_loc" | tee -a "$LOG_FILE"
                    echo "---" >> "$LOG_FILE"
                fi
            done
        fi
    done
done

echo -e "\033[32mV\033[0m Validation complete. See $LOG_FILE for any mismatches."
