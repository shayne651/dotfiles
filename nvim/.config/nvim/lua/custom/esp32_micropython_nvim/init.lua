-- File: ~/.config/nvim/lua/plugins/esp32_micropython_nvim.lua

local M = {}

M.selected_device = nil

function M.wipe_esp32()
  if M.selected_device ~= nil then 
    print(M.selected_device)
  end
end

function M.flash_esp32()
  if not M.selected_device then
    print("No device selected. Run :lua require('your_module_name').choose_device() first.")
    return
  end

  if not M.selected_version then
    print("No MicroPython firmware selected.")
    M.list_and_download_micropython_version()
    return
  end

  local firmware_path = vim.fn.stdpath('data') .. "/micropython/" .. M.selected_version
  if vim.fn.filereadable(firmware_path) == 0 then
    print("Firmware not found. Downloading...")
    M.download_micropython(M.selected_version)
    return
  end

  print("Flashing " .. M.selected_version .. " to " .. M.selected_device)

  -- Construct the esptool.py command
  local command = string.format(
    "esptool.py --chip esp32 --port %s erase_flash && esptool.py --chip esp32 --port %s write_flash -z 0x1000 %s",
    M.selected_device,
    M.selected_device,
    firmware_path
  )

  -- Run the flash command
  vim.fn.system(command)

  if vim.v.shell_error == 0 then
    print("Successfully flashed " .. M.selected_version)
  else
    print("Failed to flash. Check the connection and try again.")
  end
end

function M.full_flash_flow()
  if not M.selected_device then
    M.choose_device()
  end

  vim.defer_fn(function()
    if not M.selected_version then
      M.list_and_download_micropython_version()
    end

    vim.defer_fn(function()
      M.flash_esp32()
    end, 2000)
  end, 2000)
end



function M.choose_device()

  local device_list = {}

  if vim.loop.os_uname().sysname == "Darwin" then
    local handle = io.popen("ls /dev/tty.* | grep -i usb")
    if handle then
      for line in handle:lines() do
        table.insert(device_list, line)
      end
      handle:close()
    end
  end

  local actions = require("telescope.actions")

  require('telescope.pickers').new({}, {
    prompt_title = "Choose a Bluetooth device",
    finder = require('telescope.finders').new_table({
      results = device_list,
    }),
    sorter = require('telescope.sorters').get_generic_fuzzy_sorter(),
    attach_mappings = function(prompt_bufnr, map)
      map("i", "<CR>", function()
        local selection = require('telescope.actions.state').get_selected_entry()
        M.selected_device = selection[1]
        actions.close(prompt_bufnr)
      end)
      return true
    end,
  }):find()
end

function M.push_current_file()
  print("Pushing current file")
end

function M.run_current()
  print("Running current file")
end

function M.push_directory_files()
  print("Pushing all files in directory")
end

function M.list_and_download_micropython_version()
  local versions = {
    "esp32-20230619-v1.20.0.bin",
    "esp32-20230510-v1.19.1.bin",
    "esp32-20230415-v1.18.0.bin",
    "esp32-20230325-v1.17.0.bin",
  }

  -- Telescope picker to select the version
  local actions = require("telescope.actions")

  require('telescope.pickers').new({}, {
    prompt_title = "Choose a MicroPython Version",
    finder = require('telescope.finders').new_table({
      results = versions,
    }),
    sorter = require('telescope.sorters').get_generic_fuzzy_sorter(),
    attach_mappings = function(prompt_bufnr, map)
      map("i", "<CR>", function()
        local selection = require('telescope.actions.state').get_selected_entry()
        M.selected_version = selection[1]
        actions.close(prompt_bufnr)
        
        -- Now download the selected version
        M.download_micropython(M.selected_version)
      end)
      return true
    end,
  }):find()
end

function M.download_micropython(version)
  if not version then
    print("No version selected!")
    return
  end

  print("Downloading version: " .. version)

  -- Build the correct download URL
  local url = "https://micropython.org/resources/firmware/" .. version

  -- Destination path
  local destination = vim.fn.stdpath('data') .. "/micropython/" .. version

  -- Ensure the directory exists
  vim.fn.mkdir(vim.fn.fnamemodify(destination, ":h"), "p")

  -- Use curl to download the file
  local curl_command = string.format("curl -L %s -o %s", url, destination)
  vim.fn.system(curl_command)

  if vim.v.shell_error == 0 then
    print("Successfully downloaded " .. version)
  else
    print("Failed to download " .. version)
  end
end


return M