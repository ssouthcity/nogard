ifdef NOGARD_GUILD
	COMMANDS_ENDPOINT := "https://discord.com/api/v10/applications/${NOGARD_APPLICATION_ID}/guilds/${NOGARD_GUILD}/commands"
else
	COMMANDS_ENDPOINT := "https://discord.com/api/v10/applications/${NOGARD_APPLICATION_ID}/commands"
endif

all: commands

commands:
	@curl -X PUT \
		-H "Content-Type: application/json" \
		-H "Authorization: Bot ${NOGARD_TOKEN}" \
		-d @commands.json \
		${COMMANDS_ENDPOINT}
