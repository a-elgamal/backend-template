import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:http/http.dart' as http;
import 'package:portal/src/apps/app.dart';
import 'package:portal/src/apps/app_service.dart';
import 'package:portal/src/stored/stored.dart';

import 'app_service_test.mocks.dart';

@GenerateNiceMocks([MockSpec<http.Client>()])
void main() {
  late MockClient mockClient;
  late AppService appService;
  setUp(() {
    mockClient = MockClient();
    appService = AppService(mockClient);
  });

  const String jsonResponse = """
        {
          "id": "1", 
          "createdBy": "user", 
          "createdAt": "2020-01-01T00:00:00.000Z", 
          "modifiedBy": "user", 
          "modifiedAt": "2020-01-02T00:00:00.000Z", 
          "content": {
            "apiKey": "abc123", 
            "disabled": false
          }}
    """;

  group("listApps should", () {
    test('return a list of undeleted apps by default', () async {
      when(mockClient.get(any)).thenAnswer((_) async => http.Response(
            '[$jsonResponse]',
            200,
          ));

      final result = await appService.listApps();
      expect(result, isA<List<Stored<App>>>());
      expect(result.length, 1);
      expect(result.first.id, '1');
      verify(mockClient.get(argThat(predicate<Uri>((u) {
        final disabledQueryParam = u.queryParameters["disabled"];
        if (disabledQueryParam != null) {
          final value = bool.tryParse(disabledQueryParam);
          if (value != null) {
            return value == false;
          }
        }
        return false;
      })))).called(1);
    });

    test('return a list of deleted apps if asked', () async {
      when(mockClient.get(any)).thenAnswer((_) async => http.Response(
            '[$jsonResponse]',
            200,
          ));

      final result = await appService.listApps(showDeleted: true);
      expect(result, isA<List<Stored<App>>>());
      expect(result.length, 1);
      expect(result.first.id, '1');
      verify(mockClient.get(
          argThat(predicate<Uri>((u) => u.queryParameters.isEmpty)))).called(1);
    });

    test('handle and throw errors on non-200 responses', () async {
      when(mockClient.get(any)).thenAnswer((_) async => http.Response(
            '{"error": {"message": "Something went wrong"}}',
            400,
          ));

      expect(appService.listApps(), throwsException);
    });
  });

  test('addApp should successfully create a new app', () async {
    when(mockClient.post(any, body: anyNamed('body')))
        .thenAnswer((_) async => http.Response(
              jsonResponse,
              200,
            ));

    final result = await appService.addApp("NewApp");
    expect(result.id, '1');
    expect(result.content.apiKey, 'abc123');
  });

  test('addApp should handle and throw errors on failure', () async {
    when(mockClient.post(any, body: anyNamed('body')))
        .thenAnswer((_) async => http.Response(
              '{"error": {"message": "Unable to create app"}}',
              400,
            ));

    expect(appService.addApp("NewApp"), throwsException);
  });

  test('getApp should return app details for a specific app', () async {
    when(mockClient.get(any)).thenAnswer((_) async => http.Response(
          jsonResponse,
          200,
        ));

    final result = await appService.getApp("1");
    expect(result.id, '1');
    expect(result.content.apiKey, 'abc123');
  });

  test('getApp should handle and throw errors on non-200 responses', () async {
    when(mockClient.get(any)).thenAnswer((_) async => http.Response(
          '{"error": {"message": "App not found"}}',
          404,
        ));

    expect(appService.getApp("1"), throwsException);
  });

  test('setDisabled should handle API responses correctly', () async {
    when(mockClient.patch(any, body: anyNamed('body')))
        .thenAnswer((_) async => http.Response(
              '',
              200,
            ));

    await appService.setDisabled("1", true);
    verify(mockClient.patch(any, body: anyNamed('body'))).called(1);
  });

  test('setDisabled should throw errors on non-200 responses', () async {
    when(mockClient.patch(any, body: anyNamed('body')))
        .thenAnswer((_) async => http.Response(
              '{"error": {"message": "Invalid request"}}',
              400,
            ));

    expect(appService.setDisabled("1", true), throwsException);
  });

  test('resetAPIKey should return new API key on success', () async {
    when(mockClient.post(any)).thenAnswer((_) async => http.Response(
          'newApiKey123',
          200,
        ));

    final result = await appService.resetAPIKey("1");
    expect(result, equals("newApiKey123"));
  });

  test('resetAPIKey should handle and throw errors on failure', () async {
    when(mockClient.post(any)).thenAnswer((_) async => http.Response(
          '{"error": {"message": "Unable to reset API key"}}',
          400,
        ));

    expect(appService.resetAPIKey("1"), throwsException);
  });
}
