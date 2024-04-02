# steam-api-service
A microservice that allows for resolving a Steam HEX from profile ID or vanity URL.

## Available Endpoints
* `/get_user_hex` - The main endpoint to query for a user's Steam HEX.  
    * `?vanity_url=xxx` - query based on user's vanity URL.  
    * `?profile_id=xxx` - query based on user's profile ID.  
    * `?query=xxx` - query based on SteamCommunity profile URL.  