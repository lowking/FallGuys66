#!/bin/bash
fyne bundle -package data -name LogoWhite logo-white.png > ../bundled.go
fyne bundle -append -package data -name LogoBlack logo-black.png >> ../bundled.go
#fyne bundle -append -package data -name FontMonaco font/LigaMonacoforPowerline.otf >> ../bundled.go
fyne bundle -append -package data -name FontSmileySansOblique font/SmileySans-Oblique.ttf >> ../bundled.go
