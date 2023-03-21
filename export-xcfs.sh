#!/bin/sh

# Adapted from https://stackoverflow.com/a/5846727

gimp -n -i -b - <<EOF
  (let*
    (
        (file's (list
          (list "xcf/combat-tileset.xcf" "combat-terrain/combat-tileset.png")
          (list "xcf/tree-v1.xcf" "combat-terrain/tree-v1.png")
        ))
        (xcf "")
        (png "")
        (image 0)
        (layer 0)
    )

    (while (pair? file's)
      (set! xcf (car (car file's)))
      (set! png (car (cdr (car file's))))

      (set! image (car (gimp-file-load RUN-NONINTERACTIVE xcf xcf)))
      (set! layer (car (gimp-image-merge-visible-layers image CLIP-TO-IMAGE)))
      (gimp-file-save RUN-NONINTERACTIVE image layer png png)
      (gimp-image-delete image)

      (set! file's (cdr file's))
    )
    (gimp-quit 0)
  )
EOF
