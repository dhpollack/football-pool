#!/bin/bash

# Create an empty array
echo "[]" >backend/seed/espn-season.json

for i in {1..18}; do
  echo "Fetching week $i..."
  # Fetch the data and extract the events
  curl -s "https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard?dates=2025&seasontype=2&week=$i" | jq '.events' >"backend/seed/week_$i.json"
done

# Combine all the files
jq -s 'add' backend/seed/week_*.json >backend/seed/espn-season.json

# Clean up the temporary files
rm backend/seed/week_*.json

echo "Done."
