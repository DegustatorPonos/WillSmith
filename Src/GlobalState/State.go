package globalstate

var State = AppState {
	VersionName: "0.5.1a",

	ChannelLengths: ChannelsParams{
		ControlChannel: 3,
		RequestChannel: 2,
		ConnectionBuffer: 1,
		LogChannel: 2,
		DownloadBuffer: 16,
	},
}

type AppState struct {
	VersionName string
	ChannelLengths ChannelsParams
}

type ChannelsParams struct {
	ControlChannel int
	RequestChannel int
	ConnectionBuffer int
	LogChannel int
	DownloadBuffer int
}
