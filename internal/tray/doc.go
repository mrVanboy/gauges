// Package tray holds the whole logic for the tray app
//
// The structure of the tray menu is:
//  |-- [x] Start on login
//  |-- Device
//      |-- [x] SERIAL | /dev/tty.xxx
//  |-- Refresh rate
//      |-- 0.1s
//      |-- 0.2s
//      |-- 0.5s
//      |-- 1s
//      |-- 2s
//      |-- 5s
//  |-- --------
//  |-- Show log
//  |-- --------
//  |-- Quit
package tray
