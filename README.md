# Description

This is a project inspired by [evediscordkm](https://github.com/xSetech/evediscordkm).
It filters killmail reports generated by redisq.zkillboard.com and delivers them to discord

# Configuration file

```
zkillboard:
  replay_from: tmp/debug #optional replays saved kill reports from the disk

channels:
  - type: debug
    constraints:
      alliances:
        - The Eclipse.
        - Goonswarm Federation
      characters:
        - Audrey Horn
  - type: record
    config:
      path: tmp/debug
  - type: discord
    config:
      webhook: https://...
    constraints:
      alliances:
        - The Eclipse.
      corporations:
        - "New Eden's Best."
      systems:
        - NU4-2G
      type: win
```


# Running
```
go get github.com/vvlad/evediscordkm
go build -o evekm github.com/vvlad/evediscordkm/cmd
./evekm -config config.yaml
```

