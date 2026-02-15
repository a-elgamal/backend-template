import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:portal/src/apps/app.dart';
import 'package:portal/src/apps/app_service.dart';
import 'package:portal/src/apps/app_view.dart';
import 'package:portal/src/apps/new_app_view.dart';
import 'package:portal/src/routes.dart';
import 'package:portal/src/stored/stored.dart';
import 'package:provider/provider.dart';
import 'app_list_view_test.mocks.dart';

@GenerateNiceMocks([MockSpec<AppService>()])
void main() {
  late MockAppService mockAppService;

  List<Stored<App>> mockApps = List<Stored<App>>.unmodifiable([
    Stored<App>(
      id: "app1",
      createdBy: "admin",
      createdAt: DateTime.now(),
      modifiedBy: "admin",
      modifiedAt: DateTime.now(),
      content: const App(apiKey: "123"),
    ),
    Stored<App>(
      id: "app2",
      createdBy: "admin",
      createdAt: DateTime.now(),
      modifiedBy: "admin",
      modifiedAt: DateTime.now(),
      content: const App(apiKey: "123"),
    ),
  ]);

  const mockErrorText = "mock error";

  setUp(() {
    mockAppService = MockAppService();
  });

  Widget createTestableWidget() {
    return Provider<AppService>(
      create: (context) => mockAppService,
      child: MaterialApp.router(
        routerConfig: GoRouter(
            initialLocation: "/apps",
            debugLogDiagnostics: true,
            routes: $appRoutes),
      ),
    );
  }

  group("Upon launching", () {
    testWidgets("Fetches list of undeleted apps", (WidgetTester tester) async {
      final completer = Completer();
      when(mockAppService.listApps(showDeleted: anyNamed('showDeleted')))
          .thenAnswer((_) => completer.future.then((_) => mockApps));

      await tester.pumpWidget(createTestableWidget());
      expect(find.byType(CircularProgressIndicator), findsOneWidget);

      completer.complete();
      await tester.pump();
      expect(find.byType(CircularProgressIndicator), findsNothing);
      expect(find.byType(ListTile), findsNWidgets(mockApps.length));
      for (final app in mockApps) {
        expect(find.text(app.id), findsOne);
      }
      verify(mockAppService.listApps(showDeleted: false)).called(1);
    });

    testWidgets('shows error dialog if an error occurs',
        (WidgetTester tester) async {
      when(mockAppService.listApps(showDeleted: anyNamed('showDeleted')))
          .thenThrow(Exception(mockErrorText));

      await tester.pumpWidget(createTestableWidget());
      await tester.pump();

      // Error dialog should appeaer
      expect(find.byType(AlertDialog), findsOne);
      expect(find.textContaining(mockErrorText), findsAny);
    });

    testWidgets("shows empty state if no apps are found",
        (WidgetTester tester) async {
      when(mockAppService.listApps(showDeleted: anyNamed('showDeleted')))
          .thenAnswer((_) async => List<Stored<App>>.empty());

      await tester.pumpWidget(createTestableWidget());
      await tester.pump();

      expect(find.byType(CircularProgressIndicator), findsNothing);
      expect(find.byType(ListTile), findsNothing);
      expect(find.byIcon(Icons.visibility), findsOne);
      verify(mockAppService.listApps(showDeleted: false)).called(1);

      // Tapping on show show deleted should call the service again
      await tester.tap(find.byIcon(Icons.visibility));
      await tester.pump();

      expect(find.byType(CircularProgressIndicator), findsNothing);
      expect(find.byType(ListTile), findsNothing);
      expect(find.byIcon(Icons.visibility), findsNothing);
      verify(mockAppService.listApps(showDeleted: true)).called(1);
    });
  });

  testWidgets("Refresh icon, refetches the list of apps",
      (WidgetTester tester) async {
    var firstCall = true;
    when(mockAppService.listApps(showDeleted: anyNamed('showDeleted')))
        .thenAnswer((_) async => firstCall ? mockApps : mockApps.sublist(0, 1));

    await tester.pumpWidget(createTestableWidget());
    await tester.pump();

    expect(find.byType(ListTile), findsNWidgets(mockApps.length));

    firstCall = false;
    await tester.tap(find.byIcon(Icons.refresh));
    await tester.pump();
    expect(find.byType(ListTile), findsOneWidget);

    // 2 calls, one on first launch and then one for refreshing the list.
    verify(mockAppService.listApps(showDeleted: false)).called(2);
  });

  group("Adding a new app", () {
    const newAppName = "app";
    Future<void> saveNewApp(WidgetTester tester) async {
      await tester.pumpWidget(createTestableWidget());
      await tester.pump();

      await tester.tap(find.byType(FloatingActionButton));
      await tester.pumpAndSettle();

      expect(find.byType(NewAppView), findsOne);

      await tester.enterText(find.byType(TextField), newAppName);
      await tester.pump();
      await tester.tap(find.text("Save"));
      await tester.pumpAndSettle();
    }

    testWidgets("successfully adds the new app to the list",
        (WidgetTester tester) async {
      when(mockAppService.listApps(showDeleted: anyNamed('showDeleted')))
          .thenAnswer((_) async => <Stored<App>>[]);

      when(mockAppService.addApp(any)).thenAnswer(
        (_) async => Stored<App>(
          id: newAppName,
          createdBy: "admin",
          createdAt: DateTime.now(),
          modifiedBy: "admin",
          modifiedAt: DateTime.now(),
          content: const App(apiKey: "123"),
        ),
      );

      await saveNewApp(tester);

      // The detail of the new app is shown
      expect(find.byType(AppView), findsOne);
      expect(find.text(newAppName), findsOne);

      verify(mockAppService.addApp(newAppName)).called(1);
    });

    testWidgets("shows an error upon failure", (WidgetTester tester) async {
      const newAppName = "app";
      when(mockAppService.listApps(showDeleted: anyNamed('showDeleted')))
          .thenAnswer((_) async => <Stored<App>>[]);

      when(mockAppService.addApp(any)).thenThrow(Exception(mockErrorText));

      await saveNewApp(tester);

      // Error dialog should appeaer
      expect(find.byType(AlertDialog), findsOne);
      expect(find.textContaining(mockErrorText), findsAny);

      verify(mockAppService.addApp(newAppName)).called(1);
    });
  });

  group("An App list tile", () {
    Future<void> openFirstAppDetail(WidgetTester tester) async {
      expect(find.byType(AppView), findsOne);
      expect(find.text(mockApps[0].id), findsOne);

      await tester.tap(find.byIcon(Icons.delete));
      await tester.pumpAndSettle();

      // Tap "yes" on the confirmation dialog
      await tester.tap(find.text("Yes"));
      await tester.pumpAndSettle();
    }

    group("is correctly updated with the latest data from app detail view", () {
      setUp(() {
        when(mockAppService.listApps(showDeleted: anyNamed('showDeleted')))
            .thenAnswer((_) async => List<Stored<App>>.from(mockApps));
        when(mockAppService.getApp(mockApps[0].id)).thenAnswer((_) async =>
            Stored<App>.copy(mockApps[0],
                content: App.copy(mockApps[0].content,
                    disabled: !mockApps[0].content.disabled)));
        when(mockAppService.setDisabled(
                mockApps[0].id, !mockApps[0].content.disabled))
            .thenAnswer((_) async {});
      });

      testWidgets("Removes deleted elements when list doesn't show deleted",
          (tester) async {
        await tester.pumpWidget(createTestableWidget());
        await tester.pump();

        expect(find.text(mockApps[1].id), findsOne);
        await tester.tap(find.text(mockApps[0].id));
        await tester.pumpAndSettle();

        expect(find.byType(AppView), findsOne);
        expect(find.text(mockApps[0].id), findsOne);

        await openFirstAppDetail(tester);

        // Go back to the list view
        await tester.tap(find.byType(BackButton));
        await tester.pumpAndSettle();

        expect(find.text(mockApps[0].id), findsNothing);
        expect(find.text(mockApps[1].id), findsOne);
      });

      testWidgets("Updates status of elements when the list shows deleted",
          (tester) async {
        await tester.pumpWidget(createTestableWidget());
        await tester.pump();

        await tester.tap(find.byIcon(Icons.visibility));
        await tester.pump();

        expect(find.byType(ListTile), findsNWidgets(mockApps.length));
        await tester.tap(find.text(mockApps[0].id));
        await tester.pumpAndSettle();

        await openFirstAppDetail(tester);

        // Go back to the list view
        await tester.tap(find.byType(BackButton));
        await tester.pumpAndSettle();

        expect(find.byType(ListTile), findsNWidgets(mockApps.length));
        expect(
            tester.widget<Text>(find.text(mockApps[0].id)).style?.color,
            isNot(equals(
                tester.widget<Text>(find.text(mockApps[1].id)).style?.color)));
      });
    });
  });
}
