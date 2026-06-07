return {
  "3rd/image.nvim",
  build = "luarocks --local install magick",
  opts = {
    backend = "kitty",
    integrations = {
      neotree = {
        enabled = true,
        clear_in_insert_mode = false,
        download_remote_images = true,
        only_render_image_at_cursor = false,
        filetypes = { "png", "jpg", "gif", "webp", "avif" },
      },
    },
    max_width_window_percentage = math.huge,
    max_height_window_percentage = math.huge,
    tmux_show_only_in_active_window = true,
    hijack_file_patterns = { "*.png", "*.jpg", "*.gif", "*.webp", "*.avif" },
  },
}
