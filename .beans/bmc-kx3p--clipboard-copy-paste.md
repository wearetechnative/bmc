---
# bmc-kx3p
title: clipboard-copy-paste
status: todo
type: task
priority: normal
created_at: 2026-05-06T20:18:51Z
updated_at: 2026-05-06T20:18:51Z
---

De oude bash-config had aparte `clipboardCopyCommand` en `clipboardPasteCommand` arrays. In de Go-rewrite is dit teruggebracht naar één `clipboard_command` string (alleen copy), en is de paste-functionaliteit weggevallen.

Voeg in de config struct aparte `copy_command` en `paste_command` velden toe (of herstel de oorspronkelijke namen), zodat gebruikers zowel het kopieer- als plak-commando kunnen configureren. De paste-functionaliteit werd in de oude versie gebruikt om de MFA-code automatisch in te vullen (bijv. in een browservenster).
