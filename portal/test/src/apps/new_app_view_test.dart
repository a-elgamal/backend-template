import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:portal/src/apps/new_app_view.dart';

void main() {
  group(
    "A New App View should",
    () {
      initNewAppView(WidgetTester widgetTester, Function(String) onSave) async {
        await widgetTester.pumpWidget(
          MaterialApp(
            home: Scaffold(
              body: NewAppView(onSave: onSave),
            ),
          ),
        );
      }

      testWidgets("saves app name", (widgetTester) async {
        var savedName = "";
        await initNewAppView(widgetTester, (name) => savedName = name);
        const expectedName = "abc";
        await widgetTester.enterText(find.byType(TextField), expectedName);
        await widgetTester.pump();
        await widgetTester.tap(find.byType(ElevatedButton));
        expect(savedName, equals(expectedName));
      });

      testWidgets("doesn't allow invalid characters", (widgetTester) async {
        var savedName = "";
        await initNewAppView(widgetTester, (name) => savedName = name);
        await widgetTester.enterText(find.byType(TextField), "Ab!2");
        await widgetTester.pump();
        await widgetTester.tap(find.byType(ElevatedButton));
        expect(savedName, equals("b2"));
      });

      testWidgets("doesn't allow more than 36 characters",
          (widgetTester) async {
        var savedName = "";
        await initNewAppView(widgetTester, (name) => savedName = name);
        expect(find.text("0/36"), findsOne);
        const expected = "012345678901234567890123456789012345";
        await widgetTester.enterText(find.byType(TextField), "${expected}6789");
        await widgetTester.pump();
        expect(find.text("36/36"), findsOne);
        await widgetTester.testTextInput.receiveAction(TextInputAction.done);
        expect(savedName, equals(expected));
      });

      testWidgets("save button is disabled until a character is entered",
          (widgetTester) async {
        await initNewAppView(widgetTester, (_) {});
        // Verify that the Save button is initially disabled.
        final saveButtonFinder = find.widgetWithText(ElevatedButton, 'Save');

        expect(widgetTester.widget<ElevatedButton>(saveButtonFinder).onPressed,
            isNull);
        await widgetTester.enterText(find.byType(TextField), "a");
        await widgetTester.pump();
        expect(widgetTester.widget<ElevatedButton>(saveButtonFinder).onPressed,
            isNotNull);
      });
    },
  );
}
