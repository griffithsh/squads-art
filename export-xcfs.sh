#!/bin/sh

# Adapted from https://stackoverflow.com/a/5846727

gimp -n -i -b - <<EOF
  (let*
    (
        (file's (list
          (list "xcf/combat-tileset.xcf" "combat-terrain/combat-tileset.png")
          (list "xcf/tree-v1.xcf" "combat-terrain/tree-v1.png")
          (list "xcf/overworld-grass-encroaching-sand.xcf" "packed/content/overworld/tiles/grass-encroaching-sand.png")
          (list "xcf/overworld-water-encroaching-sand.xcf" "packed/content/overworld/tiles/water-encroaching-sand-edges.png")
          (list "xcf/overworld-sand-encroaching-water.xcf" "packed/content/overworld/tiles/sand-encroaching-water-edges.png")
          (list "xcf/overworld-water-encroaching-sand-corners.xcf" "packed/content/overworld/tiles/water-encroaching-sand-corners.png")
          (list "xcf/overworld-sand-encroaching-water-corners.xcf" "packed/content/overworld/tiles/sand-encroaching-water-corners.png")
          (list "xcf/overworld-grass-base.xcf" "packed/content/overworld/tiles/grass.png")
          (list "xcf/overworld-sand-base.xcf" "packed/content/overworld/tiles/sand.png")
          (list "xcf/overworld-water-base.xcf" "packed/content/overworld/tiles/water.png")
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
