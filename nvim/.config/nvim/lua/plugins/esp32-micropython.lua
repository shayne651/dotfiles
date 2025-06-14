return {
    dir = vim.fn.stdpath("config") .. "/lua/custom/esp32_micropython_nvim",
    name = "esp32_micropython_nvim",
    lazy = false,
    config = function()
      local esp32 = require("custom.esp32_micropython_nvim")
  
      vim.api.nvim_create_user_command("WipeEsp32", esp32.wipe_esp32, {})
      vim.api.nvim_create_user_command("FlashEsp32", esp32.flash_esp32, {})
      vim.api.nvim_create_user_command("RunEsp32", esp32.run_current, {})
      vim.api.nvim_create_user_command("PushFileEsp32", esp32.push_current_file, {})
      vim.api.nvim_create_user_command("DownloadMicroPython", esp32.list_and_download_micropython_version, {})
      vim.keymap.set('n', '<leader>dm', esp32.list_and_download_micropython_version, { desc = 'Download MicroPython' })
      vim.keymap.set('n', '<leader>b', esp32.choose_device, { desc = 'ESP32: List Devices' })
      vim.keymap.set('n', '<leader>w', esp32.wipe_esp32, { desc = 'ESP32: Wipe' })
      vim.keymap.set('n', '<leader>r', esp32.run_current, { desc = 'ESP32: Run File' })
    end,
  }