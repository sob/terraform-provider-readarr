resource "readarr_download_client_utorrent" "example" {
  enable        = true
  priority      = 1
  name          = "Example"
  host          = "utorrent"
  url_base      = "/utorrent/"
  port          = 9091
  book_category = "tv-readarr"
}