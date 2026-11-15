from datetime import datetime, timezone
from PIL import Image, ImageDraw

# Make a simple test picture.
img = Image.new("RGB", (320, 200), "white")
draw = ImageDraw.Draw(img)
draw.rectangle((20, 20, 150, 180), fill="cornflowerblue")
draw.ellipse((170, 35, 295, 160), fill="orange")
draw.text((20, 170), "EXIF test image", fill="black")

# Use one timestamp for the EXIF date fields.
now_local = datetime.now().astimezone()
exif_datetime = now_local.strftime("%Y:%m:%d %H:%M:%S")
creation_date_utc = datetime.now(timezone.utc).isoformat(timespec="seconds")

# Sample EXIF values—no personal or location data.
exif = Image.Exif()
exif[271] = "Example Camera"  # Make
exif[272] = "Test Model"      # Model
exif[306] = exif_datetime     # DateTime

# DateTimeOriginal and DateTimeDigitized belong in the EXIF sub-IFD.
exif_ifd = {
    36867: exif_datetime,  # DateTimeOriginal
    36868: exif_datetime,  # DateTimeDigitized
}
exif[34665] = exif_ifd      # Exif IFD pointer

exif_bytes = exif.tobytes()

# EXIF-bearing formats
img.save("sample.jpeg", "JPEG", exif=exif_bytes)
img.save("sample.png", "PNG", exif=exif_bytes)
img.save("sample.tiff", "TIFF", exif=exif_bytes)

# GIF: custom comment fields, not EXIF fields.
gif_comment = (
    f"DateTimeOriginal={exif_datetime}; "
    f"CreationDate={creation_date_utc}"
).encode("utf-8")
img.save("sample.gif", "GIF", comment=gif_comment)

# BMP: image only; no standard EXIF added.
img.save("sample.bmp", "BMP")

print(
    "Created sample.jpeg, sample.png, sample.tiff, "
    "sample.gif, and sample.bmp"
)
