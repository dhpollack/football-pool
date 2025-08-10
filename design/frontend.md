# Frontend Design

## Libraries

We want a react frontend for our app.  The site should require a login to use.  There will be two types of users, players and administrators.  The players should have limited permissions to interact with the player pages and the administrators should have all the player permissions plus special admin permissions to do administrative tasks for the site.

## Pages

- Homepage: landing page of the site with a navigation drawer on the left and a login/profile at the top right
- Pick Entry / Pick View: Players will enter their weekly picks for the weeks game.  There will be an favorite and an underdog and a spread for each game.  The player must pick either the favorite or the underdog and assign a rank value for the pick between 1 and N, where N is the total number of games for that week.  The rank for each game must be unique.  These entries should go into a database table with the game id, the player id, and the pick and the rank value. There should also be a "Quick Pick" option that selects a favorite or underdog at random and assigns a random rank.  The "Quick Pick" option needs to conform to the constraints of the entry, i.e. rank values must be unique.
- Result Entry: (admin) Administrators should be able to enter the scores for each game.  These should go into a database table with the game id, the favorite's score, the underdog's score and a result that is either the favorite wins (favorite - spread > underdog), the underdog wins (underdog > favorite - spread), or a tie (favorite - spread = underdog).
- Week's Results: We should be able to view the results for all players for a single week.  This should show a table of the sum of the rank of the correctly selected results of the games of the players' in descending order.  The player should receive the rank value they assigned for correctly picked results and 0 points for incorrect picks.  For ties, all players should receive half the amount of points of the rank that they assigned to the game.  There should be an icon to indicate if the player used the "Quick Pick" option to make their entry.
- Overall Results: Similar to the Week's Results, we want to see the overall results.  This should be the sum of the weekly results for a single season for all players.  The order of the results should be in descending order.
- Login: If site visitor is not logged in, then they should be redirected to a login page
- Player Profile: There should be a page for the player to add profile information such as name and home address
