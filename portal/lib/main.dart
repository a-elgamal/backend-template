// ignore_for_file: avoid_print

import 'package:flutter/material.dart';
import 'package:portal/constants.dart';

import 'src/app.dart';

void main() async {
  print("backend scheme: $kBackendScheme");
  print("backend host: $kBackendHost");
  print("backend port: $kBackendPort");
  print("backend path: $kBackendPath");
  print("backend path segments: $kBackendSegments");

  // Run the app and pass in the SettingsController. The app listens to the
  // SettingsController for changes, then passes it further down to the
  // SettingsView.
  runApp(const MyApp());
}
