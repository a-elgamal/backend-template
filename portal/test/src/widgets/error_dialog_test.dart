import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:portal/src/widgets/error_dialog.dart';

void main() {
  testWidgets(
      'showError displays an error dialog with the correct message and error details',
      (WidgetTester tester) async {
    // Provide a context within which to display the dialog
    await tester.pumpWidget(MaterialApp(
      home: Scaffold(
        body: Builder(
          builder: (BuildContext context) {
            return ElevatedButton(
              onPressed: () {
                showError(context, 'Sample Error', StackTrace.current,
                    'An unexpected error occurred.');
              },
              child: const Text('Trigger Error'),
            );
          },
        ),
      ),
    ));

    // Simulate button press to trigger the error dialog
    await tester.tap(find.text('Trigger Error'));
    await tester.pumpAndSettle(); // Wait for the dialog to appear

    // Verify the dialog title, message, and error details are displayed correctly
    expect(find.text('Error'), findsOneWidget);
    expect(find.text('An unexpected error occurred.'), findsOneWidget);
    expect(find.text('Sample Error'), findsOneWidget);

    // Check that the Close button is present and functional
    expect(find.text('Close'), findsOneWidget);
    await tester.tap(find.text('Close'));
    await tester.pumpAndSettle(); // Dismiss the dialog

    // Ensure the dialog has been dismissed
    expect(find.byType(AlertDialog), findsNothing);
  });
}
