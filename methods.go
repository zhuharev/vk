package vk

const (
	METHOD_ACCOUNT_GET_BANNED  = "account.getBanned"
	METHOD_ACCOUNT_SET_ONLINE  = "account.setOnline"
	METHOD_ACCOUNT_SET_OFFLINE = "account.setOffline"

	METHOD_AUDIO_GET           = "audio.get"
	METHOD_AUDIO_SET_BROADCAST = "audio.setBroadcast"

	METHOD_AUTH_SIGNUP  = "auth.signup"
	METHOD_AUTH_CONFIRM = "auth.confirm"

	METHOD_BOARD_GET_COMMENTS = "board.getComments"

	METHOD_GROUPS_GET         = "groups.get"
	METHOD_GROUPS_GET_BY_ID   = "groups.getById"
	METHOD_GROUPS_GET_MEMBERS = "groups.getMembers"
	METHOD_GROUPS_JOIN        = "groups.join"
	METHOD_GROUPS_LEAVE       = "groups.leave"
	METHOD_GROUPS_SEARCH      = "groups.search"

	METHOD_MESSAGES_SEND         = "messages.send"
	METHOD_MESSAGES_GET          = "messages.get"
	METHOD_MESSAGES_GET_DIALOGS  = "messages.getDialogs"
	METHOD_MESSAGES_MARK_AS_READ = "messages.markAsRead"
	METHOD_MESSAGES_GET_HISTORY  = "messages.getHistory"
	METHOD_MESSAGES_SET_ACTIVITY = "messages.setActivity"

	METHOD_LIKES_ADD      = "likes.add"
	METHOD_LIKES_DELETE   = "likes.delete"
	METHOD_LIKES_IS_LIKED = "likes.isLiked"
	METHOD_LIKES_GET_LIST = "likes.getList"

	METHOD_WALL_GET         = "wall.get"
	METHOD_WALL_GET_BY_ID   = "wall.getById"
	METHOD_WALL_POST        = "wall.post"
	METHOD_WALL_REPOST      = "wall.repost"
	METHOD_WALL_DELETE      = "wall.delete"
	METHOD_WALL_GET_REPOSTS = "wall.getReposts"

	METHOD_FRIENDS_GET             = "friends.get"
	METHOD_FRIENDS_GET_REQUESTS    = "friends.getRequests"
	METHOD_FRIENDS_ARE_FRIENDS     = "friends.areFriends"
	METHOD_FRIENDS_ADD             = "friends.add"
	METHOD_FRIENDS_DELETE          = "friends.delete"
	METHOD_FRIENDS_GET_MUTUAL      = "friends.getMutual"
	METHOD_FRIENDS_GET_SUGGESTIONS = "friends.getSuggestions"

	METHOD_PHOTOS_GET                    = "photos.get"
	METHOD_PHOTOS_GET_ALL                = "photos.getAll"
	METHOD_PHOTOS_GET_WALL_UPLOAD_SERVER = "photos.getWallUploadServer"
	METHOD_PHOTOS_SAVE_WALL_PHOTO        = "photos.saveWallPhoto"

	METHOD_USERS_GET               = "users.get"
	METHOD_USERS_GET_FOLLOWERS     = "users.getFollowers"
	METHOD_USERS_SEARCH            = "users.search"
	METHOD_USERS_GET_SUBSCRIPTIONS = "users.getSubscriptions"

	METHOD_STATUS_SET = "status.set"

	METHOD_DATABASE_GET_CITIES       = "database.getCities"
	METHOD_DATABASE_GET_CITIES_BY_ID = "database.getCitiesById"

	METHOD_UTILS_RESOLVE_SCREEN_NAME = "utils.resolveScreenName"

	METHOD_NEWSFEED_ADD_BAN = "newsfeed.addBan"

	METHOD_EXECUTE = "execute"
)
