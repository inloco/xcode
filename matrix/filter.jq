{
    include: (
        to_entries | map(
            select(
                (
                    .value.version.build != "8E1000a"
                ) and (
                    .value.version.release.gm or .value.version.release.release
                ) and (
                    .value.links.download.url // "" | endswith(".xip")
                )
            ) | {
                name: .value.name,
                version: .value.version.number,
                build: .value.version.build,
                url: .value.links.download.url,
                latest: (
                    .key == 0
                ),
            }
        ) | sort_by(
            .version | split(".") | map(tonumber)
        ) | reverse
    ),
}
