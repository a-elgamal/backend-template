import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:portal/src/widgets/confirmation_dialog.dart';

void main() {
  const String title = 'Confirm Action';
  const String text = 'Are you sure you want to proceed?';
  const String yesLabel = 'Proceed';
  const String noLabel = 'Cancel';
  const String showDialogLabel = 'Trigger Confirmation';

  testWidgets(
      'showError displays an error dialog with the correct message and error details',
      (WidgetTester tester) async {
    bool? result;

    // Provide a context within which to display the dialog
    await tester.pumpWidget(MaterialApp(
      home: Scaffold(
        body: Builder(
          builder: (BuildContext context) {
            return ElevatedButton(
              onPressed: () {
                showConfirmationDialog(
                  context,
                  title: title,
                  text: text,
                  yes: yesLabel,
                  no: noLabel,
                ).then((value) => result = value);
              },
              child: const Text(showDialogLabel),
            );
          },
        ),
      ),
    ));

    // Simulate button press to trigger the confirmation dialog
    await tester.tap(find.text(showDialogLabel));
    await tester.pumpAndSettle(); // Wait for the dialog to appear

    // Verify the dialog title, message, and error details are displayed correctly
    expect(find.text(title), findsOneWidget);
    expect(find.text(text), findsOneWidget);
    expect(find.text(yesLabel), findsOneWidget);
    expect(find.text(noLabel), findsOneWidget);

    // Check that the yes button is present and functional
    await tester.tap(find.text(yesLabel));
    await tester.pumpAndSettle(); // Dismiss the dialog

    // Ensure the dialog has been dismissed and true is returned
    expect(find.byType(AlertDialog), findsNothing);
    expect(result, isTrue);

    // Retrigger the same dialog
    // Simulate button press to trigger the confirmation dialog
    await tester.tap(find.text(showDialogLabel));
    await tester.pumpAndSettle(); // Wait for the dialog to appear

    // Check that the no button is present and functional
    await tester.tap(find.text(noLabel));
    await tester.pumpAndSettle(); // Dismiss the dialog

    // Ensure the dialog has been dismissed
    expect(find.byType(AlertDialog), findsNothing);
    expect(result, isFalse);
  });
}
