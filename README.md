# RSS-Feed
Generates a feed based on subscribed RSS content

## Installation

>NOTE: This project requires an installation of postgresql and Go 

Install the RSS-Feed cli using:

`go install "github.com/NHemmerly/RSS-Feed@latest"`

RSS-Feed requires a config file located in the user's home directory named `.gatorconfig.json` and following the format:

```
{
    "db_url": "postgres://[postgres username]:[password]@localhost:5432/gator"
    "current_user_name":"[can be left blank]" 
}
```
## Usage
### Commands
#### Register
`RSS-Feed register <username>`

Registers a new user and makes that user the current user. This should be the first command you run after installing RSS-Feed cli. This command will allow the user to start following available feeds. 
#### Login
`RSS-Feed login <username>`

Makes the specified username the current user. 
#### Reset
`RSS-Feed reset`

Deletes all existing user data.
#### Users
`RSS-Feed users`

Lists all available users, marking the current user. 
#### Agg
`RSS-Feed agg <duration>`

Starts aggregating posts from the current user's followed feeds. Agg will check for new posts every cycle of the specified `duration`.
#### Add Feed
`RSS-Feed addfeed <feed url> <name>`

Adds a new RSS feed to the available feeds.
#### Feeds
`RSS-Feed feeds`

Lists the available feeds that the current user can follow.
#### Follow
`RSS-Feed follow <url>`

Current user follows the feed specified by `url` if it is found in the available feeds.
#### Unfollow
`RSS-Feed unfollow <url>`

Current user unfollows the feed specified by `url`.
#### Following
`RSS-Feed following`

Lists the names of all feeds that the current user is following.
#### Browse
`RSS-Feed browse <optional: number of posts>`

Displays the `number of posts` most recent posts followed by the current user and collected by the aggregate command. 
