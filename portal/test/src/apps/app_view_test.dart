import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:portal/src/apps/app.dart';
import 'package:portal/src/apps/app_service.dart';
import 'package:portal/src/apps/app_view.dart';
import 'package:portal/src/stored/stored.dart';
import 'app_view_test.mocks.dart';

@GenerateNiceMocks([MockSpec<AppService>()])
void main() {
  late MockAppService mockAppService;

  final testApp = Stored<App>(
    id: "app1",
    createdBy: "creator",
    createdAt: DateTime.now().add(const Duration(days: -1)),
    modifiedBy: "modifier",
    modifiedAt: DateTime.now(),
    content: const App(apiKey: "123456", disabled: false),
  );

  const mockErrorText = "mock error";

  setUp(() {
    mockAppService = MockAppService();
  });

  Widget createAppViewWidget(Stored<App>? app) {
    return MaterialApp(
      home: AppView(mockAppService, testApp.id, app: app),
    );
  }

  group("Upon opening", () {
    void expectAppDetailShown(WidgetTester tester, Stored<App> expectedApp) {
      expect(find.text(expectedApp.createdBy), findsOneWidget);
      expect(find.text(expectedApp.createdAt.toLocal().toString()),
          findsOneWidget);
      expect(find.text(expectedApp.modifiedBy), findsOneWidget);
      expect(find.text(expectedApp.modifiedAt.toLocal().toString()),
          findsOneWidget);
      expect(find.text(expectedApp.content.apiKey), findsOneWidget);
    }

    testWidgets('displayed provided app details', (WidgetTester tester) async {
      await tester.pumpWidget(createAppViewWidget(testApp));
      expect(find.byType(CircularProgressIndicator), findsNothing);

      // Simulate async load
      await tester.pump();
      expectAppDetailShown(tester, testApp);
      verifyNever(mockAppService.getApp(testApp.id));
    });
    testWidgets('fetches app details if not provided',
        (WidgetTester tester) async {
      when(mockAppService.getApp(any)).thenAnswer((_) async => testApp);

      await tester.pumpWidget(createAppViewWidget(null));
      expect(find.byType(CircularProgressIndicator), findsOneWidget);

      // Simulate async load
      await tester.pumpAndSettle();

      expectAppDetailShown(tester, testApp);
      verify(mockAppService.getApp(testApp.id)).called(1);
    });

    testWidgets(
        'Shows error dialog if failure occurs while fetching app details',
        (WidgetTester tester) async {
      when(mockAppService.getApp(any)).thenThrow(Exception(mockErrorText));

      await tester.pumpWidget(createAppViewWidget(null));
      await tester.pump(); // Load initial state

      // Simulate the error being handled and error dialog being shown
      await tester.pump();

      // Error dialog should appeaer
      expect(find.byType(AlertDialog), findsOne);
      expect(find.textContaining(mockErrorText), findsAny);
    });

    testWidgets('refreshes app details if refresh is pressed',
        (WidgetTester tester) async {
      final newApp = Stored<App>(
        id: testApp.id,
        createdBy: "admin3@kynzy.com",
        createdAt: DateTime.now().add(const Duration(days: -1)),
        modifiedBy: "admin4@kynzy.com",
        modifiedAt: DateTime.now(),
        content: const App(
          apiKey: "New Key",
          disabled: true,
        ),
      );
      when(mockAppService.getApp(any)).thenAnswer((_) async => newApp);

      await tester.pumpWidget(createAppViewWidget(testApp));

      // Simulate load
      await tester.pump();
      expectAppDetailShown(tester, testApp);

      // Tap refresh
      await tester.tap(find.byIcon(Icons.refresh));
      await tester.pump();
      expectAppDetailShown(tester, newApp);
      verify(mockAppService.getApp(testApp.id)).called(1);
    });
  });

  group("Reset API Key", () {
    Future<void> confirmResetAPIKey(WidgetTester tester, bool yes) async {
      await tester.pumpWidget(createAppViewWidget(testApp));
      await tester.pump(); // Ensure the initial state is loaded

      await tester.tap(find.byIcon(Icons.create));
      await tester.pump(); // Start the confirmation dialog

      // Mock dialog response
      await tester.tap(find.widgetWithText(TextButton, yes ? "Yes" : "No"));
      await tester.pumpAndSettle();
    }

    testWidgets('shows a confirmation dialog and resets if yes is pressed',
        (WidgetTester tester) async {
      const newAPIKey = "newKey";
      when(mockAppService.resetAPIKey(any)).thenAnswer((_) async => newAPIKey);
      when(mockAppService.getApp(testApp.id)).thenAnswer((_) async =>
          Stored<App>.copy(testApp,
              content: App.copy(testApp.content, apiKey: newAPIKey)));

      await confirmResetAPIKey(tester, true);

      verify(mockAppService.resetAPIKey(testApp.id)).called(1);
      expect(find.text(newAPIKey), findsOne);
    });

    testWidgets('shows an error if an error occurs while reseting API Key',
        (WidgetTester tester) async {
      when(mockAppService.resetAPIKey(any)).thenThrow(Exception(mockErrorText));

      await confirmResetAPIKey(tester, true);

      // Error dialog should appeaer
      expect(find.byType(AlertDialog), findsOne);
      expect(find.textContaining(mockErrorText), findsAny);

      expect(find.text(testApp.content.apiKey), findsOne);
    });

    testWidgets('does nothing if no is pressed on the confirmation dialog',
        (WidgetTester tester) async {
      await tester.pumpWidget(createAppViewWidget(testApp));
      await tester.pump(); // Ensure the initial state is loaded

      await confirmResetAPIKey(tester, false);

      verifyNever(mockAppService.resetAPIKey(testApp.id));
      expect(find.text(testApp.content.apiKey), findsOne);
    });
  });

  group("Delete/Restore app", () {
    Future<void> deleteApp(WidgetTester tester, bool confirm) async {
      await tester.pumpWidget(createAppViewWidget(testApp));
      await tester.pump(); // Load the initial state

      await tester.tap(find.byIcon(Icons.delete));
      await tester.pump(); // Start the confirmation dialog

      // Tap action on the confirmation dialog
      await tester.tap(find.widgetWithText(TextButton, confirm ? "Yes" : "No"));
      await tester.pump();
    }

    testWidgets('Toggle app state when done twice',
        (WidgetTester tester) async {
      when(mockAppService.getApp(testApp.id)).thenAnswer((_) async => testApp);
      when(mockAppService.setDisabled(testApp.id, any))
          .thenAnswer((_) async {});

      await deleteApp(tester, true);

      verify(mockAppService.setDisabled(testApp.id, true)).called(1);

      // Toggle Restoration
      await tester.tap(find.byIcon(Icons.restore));
      await tester.pump(); // Start the confirmation dialog

      // Tap "yes" on the confirmation dialog
      await tester.tap(find.widgetWithText(TextButton, "Yes"));
      await tester.pump();

      verify(mockAppService.setDisabled(testApp.id, false)).called(1);
    });

    testWidgets('does nothing if No is pressed on the confirmation dialog',
        (WidgetTester tester) async {
      await deleteApp(tester, false);

      verifyNever(mockAppService.setDisabled(testApp.id, true));
    });

    testWidgets('shows an error message if an error occurs while deleting',
        (WidgetTester tester) async {
      when(mockAppService.setDisabled(testApp.id, any))
          .thenThrow(Exception(mockErrorText));

      await deleteApp(tester, true);

      verify(mockAppService.setDisabled(testApp.id, true)).called(1);

      await tester.pumpAndSettle();

      expect(find.byType(AlertDialog), findsOne);
      expect(find.textContaining(mockErrorText), findsAny);
      expect(find.byIcon(Icons.delete), findsOne);
    });
  });

  testWidgets("Returns app details when back button is pressed",
      (WidgetTester tester) async {
    Stored<App>? result;

    await tester.pumpWidget(MaterialApp(
      home: Builder(
        builder: (context) => ElevatedButton(
            child: const Center(),
            onPressed: () async => result = await Navigator.of(context).push(
                  MaterialPageRoute(
                    builder: (_) =>
                        AppView(mockAppService, testApp.id, app: testApp),
                  ),
                )),
      ),
    ));
    await tester.pumpAndSettle();

    // Open the page
    await tester.tap(find.byType(ElevatedButton));
    await tester.pumpAndSettle();
    expect(find.byType(AppView), findsOne);

    // tap back
    await tester.tap(find.byType(BackButton));
    await tester.pumpAndSettle();
    expect(find.byType(AppView), findsNothing);
    expect(result, equals(testApp));
  });
}
