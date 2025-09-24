.[] | . as $event |
.competitions[] | . as $competition |
($competition.competitors | map(select(.homeAway == "home"))[0].team) as $home_team |
($competition.competitors | map(select(.homeAway == "away"))[0].team) as $away_team |
($competition.odds[0]) as $odds |
{
  week: $event.week.number,
  season: $event.season.year,
  favorite_team: (if $odds.homeTeamOdds.favorite then $home_team.displayName else $away_team.displayName end),
  underdog_team: (if $odds.homeTeamOdds.favorite then $away_team.displayName else $home_team.displayName end),
  spread: ($odds.spread | if . == null then 0 else . end | if . < 0 then . * -1 else . end),
  start_time: ($event.date | sub("Z"; ":00Z"; "i"))
}
