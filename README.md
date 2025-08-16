# FILIPUS DEV BR

## Tech Stack

Less is better, no front-end bloated framework crap!
Pure juice of SSR (serve side reendering) with a RESTful Hypermedia-Driven App.

- Go
- [Templ](https://templ.guide])
  - version: v0.3.898
- [Tailwind CSS](https://tailwindcss.com/docs/installation/tailwind-cli)
- [daisyUI](https://daisyui.com/docs/install/cli/)
- [htmx](https://htmx.org/)

Dev Env dependencies:

- [air](https://github.com/air-verse/air)
- [tint](https://github.com/lmittmann/tint)
- [inotify-tools](https://github.com/inotify-tools/inotify-tools)
- [richgo](https://github.com/kyoh86/richgo)

## Project Organization

## How to dev:

### TL;DR

1. `make setup`
<!-- TODO:  Check how to do it in a more elegant way
            Include into Makefile?
-->
2. `go get github.com/a-h/templ@v0.3.898`
3. `go get github.com/lmittmann/tint`
<!-- TODO:  Make it more "ubuntu" friendly
            Not everyone uses Arch Linux - what a shame!
            Include into Makefile?
-->
4. `yay -S inotify-tools` OR
   - `sudo pacman -S inotify-tools`
   - `sudo apt install inotify-tools` (never tested!)
5. `go install github.com/kyoh86/richgo@latest`
6. `make dev-up`

### 1. make setup

Will run a setup check and, if needed, install packages, libs, so on.

### 2. go get github.com/a-h/templ@v0.3.898

Will install the templ module into the project.

### 3. go get github.com/lmittmann/tint

Will install "kinda of" formatter for slog module.

### 4. yay -S inotify-tools

Will install a tool needed to "watch" the go test command.

### 5. go install github.com/kyoh86/richgo@latest

Will install a colorful output for go test command.

### 6. make dev-up

Starts all needed watchers:

- `tailwind`: to watch css changes.
- `templ`: to watch templ file changes.
- `server`: to watch .go files' changes and live reload.
- `test-watch`: to watch .go files changes and live run the tests.
