import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:portal/src/app.dart';
import 'package:uuid/v4.dart';

void apps() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  // Returns the name of the added app
  Future<String> addApp(WidgetTester tester) async {
    await tester.pumpWidget(const MyApp());

    await tester.tap(find.byType(FloatingActionButton));

    await tester.pumpAndSettle();

    // Dashes are not accepted in the app name
    final newAppName = const UuidV4().generate().replaceAll(RegExp(r'-'), "");
    await tester.enterText(find.byType(TextField), newAppName);
    await tester.pumpAndSettle(); // En
    await tester.tap(find.text("Save"));

    // Navigates to App View Page
    await tester.pumpAndSettle();

    // App ID and an API Key in the page header
    expect(find.text(newAppName), findsOne);
    expect(tester.widget<SelectableText>(find.byKey(const Key("api-key"))).data,
        isNotNull);

    return newAppName;
  }

  testWidgets('Deletes and Restores an app successfully', (tester) async {
    final appName = await addApp(tester);

    await tester.tap(find.byIcon(Icons.delete));

    // Confirm deletion
    await tester.pumpAndSettle();
    await tester.tap(find.text("Yes"));

    // Go back to the list
    await tester.pumpAndSettle();
    await tester.tap(find.byType(BackButton));

    // App shouldn't be found
    await tester.pumpAndSettle();
    expect(find.text(appName), findsNothing);

    // Show deleted, shows the app
    await tester.tap(find.byIcon(Icons.visibility));
    await tester.pumpAndSettle();

    // Reopen app detail
    await tester.dragUntilVisible(
        find.text(appName), find.byType(ListView), const Offset(0, -100));
    await tester.tap(find.text(appName));
    await tester.pumpAndSettle();

    // Restore the app
    await tester.tap(find.byIcon(Icons.restore));
    await tester.pumpAndSettle();
    await tester.tap(find.text("Yes"));
    await tester.pumpAndSettle();

    // Go back to the list and hide deleted
    await tester.pumpAndSettle();
    await tester.tap(find.byType(BackButton));
    await tester.pumpAndSettle();
    await tester.tap(find.byIcon(Icons.visibility_off));
    await tester.pumpAndSettle();

    await tester.dragUntilVisible(
        find.text(appName), find.byType(ListView), const Offset(0, -100));
  });

  testWidgets('Resets an api-key successfully', (tester) async {
    await addApp(tester);

    // Fetch current API key
    expect(find.byKey(const Key("api-key")), findsOne);
    final originalKey =
        tester.widget<SelectableText>(find.byKey(const Key("api-key"))).data;

    await tester.tap(find.byIcon(Icons.create));

    // Confirm reset
    await tester.pumpAndSettle();
    await tester.tap(find.text("Yes"));

    // The original API Key shouldn't be found
    await tester.pumpAndSettle();
    expect(find.text(originalKey!), findsNothing);
  });
}
