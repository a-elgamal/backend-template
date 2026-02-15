import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';

void showError(BuildContext context, final Object error,
    final StackTrace stackTrace, final String message) {
  if (kDebugMode) {
    print("Error: $error");
    print("StackTrace: $stackTrace");
  }
  showDialog(
    context: context,
    builder: (ctx) => AlertDialog(
      title: const Text("Error"),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(message),
          Text("$error"),
        ],
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text("Close"),
        )
      ],
    ),
  );
}
