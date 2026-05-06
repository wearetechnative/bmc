---
# bmc-p4wn
title: post-rewrite-issues
status: completed
type: task
priority: normal
created_at: 2026-05-06T10:30:00Z
updated_at: 2026-05-06T14:00:00Z
openspec-link: openspec/changes/archive/2026-05-06-post-rewrite-fixes/proposal.md
---

Bugs en ontbrekende functionaliteit gevonden na de Go-rewrite.

## Issue 1: Lege regels in TUI-keuzelijsten

`ui.Choose()` berekent de lijsthoogte als `len(items) + 6`. Voor korte lijsten (2–4 items, zoals de "Connection method"- en "SSH user"-pickers in `ec2connect`) levert dit 3–4 lege item-slots op tussen de laatste optie en de help bar.

**Oplossing**: hoogte aanpassen naar `len(items) + 3`; voor ≤ 4 items `SetShowHelp(false)` en hoogte `len(items) + 1`.

## Issue 3: `bmc profsel` geeft geen waarschuwing als shell-wrapper ontbreekt

`bmc profsel` print `export AWS_PROFILE=xxx` naar stdout maar voert de export niet uit als de shell-wrapper niet geïnstalleerd is. Er is geen foutmelding of hint — de gebruiker ziet de TUI, kiest een profiel, en er gebeurt niets zichtbaars.

**Oplossing**: detecteer of stdout een TTY is. Als stdout een terminal is (niet een pipe via `eval "$(...)`"), dan weet de binary dat de output niet gecaptured wordt en kan het op stderr een waarschuwing tonen:
> ⚠ Shell wrapper not installed. Run: bmc install-shell-integration

Als stdout een pipe is (eval-context) blijft de output schoon.

README is al correct gedocumenteerd (Setup stap 1), maar de binary zelf geeft geen feedback op het moment dat het misgaat.

## Issue 4: `bmc install-shell-integration` faalt op NixOS met `permission denied`

Op NixOS met home-manager is `~/.zshrc` een read-only symlink naar de Nix store (`/nix/store/xxx-home-manager-files/.zshrc`). `bmc install-shell-integration` probeert ernaar te schrijven met `O_APPEND|O_WRONLY` en krijgt `permission denied`. Er volgt een generieke foutmelding zonder uitleg of alternatief.

**Werkelijke foutmelding**:
```
Error: failed to open /home/wtoorren/.zshrc: open /home/wtoorren/.zshrc: permission denied
```

**Oplossing**: vang `permission denied` op bij het openen van het rc-bestand en toon een nuttige fallback met wrapper-snippets voor alle gangbare setups:
- home-manager (`programs.zsh.initContent` / `programs.bash.initContent`)
- handmatig bash/zsh
- Fish shell (`~/.config/fish/config.fish`)

Geen NixOS-detectie nodig — `permission denied` is zelf het signaal dat het bestand beheerd wordt door een extern systeem. Dit werkt ook voor andere situaties buiten NixOS.

## Issue 2: `bmc console` opent geen containerized tab

`bmc console` opent de AWS-console in de browser maar mist de containerized-tab-functionaliteit die `assumego` (Granted) bood. Granted opent de console in een geïsoleerde browsertab per AWS-profiel zodat sessies niet door elkaar lopen. Deze isolatie ontbreekt volledig in de huidige implementatie.
