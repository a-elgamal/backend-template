import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:portal/constants.dart';
import 'package:portal/src/apps/app.dart';
import 'package:portal/src/stored/stored.dart';

class AppService {
  static final relativePathSegments = [...kBackendSegments, "apps"];

  final http.Client client;

  AppService(this.client);

  Future<List<Stored<App>>> listApps({bool showDeleted = false}) async {
    var response = await client.get(
      Uri(
        scheme: kBackendScheme,
        host: kBackendHost,
        port: kBackendPort,
        pathSegments: relativePathSegments,
        queryParameters: showDeleted ? null : {App.disabledJSONKey: "false"},
      ),
    );

    if (response.statusCode != 200) {
      handleErrorResponse(response);
    }

    List<Stored<App>> result = [];
    for (final obj in json.decode(response.body)) {
      result.add(Stored.fromJSON(App.fromJSON, obj));
    }

    return result;
  }

  Future<Stored<App>> addApp(String name) async {
    var response = await client.post(
      Uri(
        scheme: kBackendScheme,
        host: kBackendHost,
        port: kBackendPort,
        pathSegments: relativePathSegments,
      ),
      body: json.encode(<String, String>{
        "id": name,
      }),
    );

    if (response.statusCode != 200) {
      handleErrorResponse(response);
    }

    return Stored.fromJSON(App.fromJSON, json.decode(response.body));
  }

  Future<Stored<App>> getApp(String name) async {
    final response = await client.get(
      Uri(
          scheme: kBackendScheme,
          host: kBackendHost,
          port: kBackendPort,
          pathSegments: [...relativePathSegments, name]),
    );

    if (response.statusCode != 200) {
      handleErrorResponse(response);
    }

    return Stored.fromJSON(App.fromJSON, json.decode(response.body));
  }

  Future<void> setDisabled(String name, bool disabled) async {
    var response = await client.patch(
      Uri(
        scheme: kBackendScheme,
        host: kBackendHost,
        port: kBackendPort,
        pathSegments: [...relativePathSegments, name],
      ),
      body: json.encode(
        <String, bool>{
          "disabled": disabled,
        },
      ),
    );

    if (response.statusCode != 200) {
      handleErrorResponse(response);
    }
  }

  void handleErrorResponse(final http.Response response) {
    if (response.body.isNotEmpty) {
      Map<String, dynamic> responseJSON = {};
      try {
        responseJSON = json.decode(response.body);
      } catch (_) {}

      if (responseJSON.containsKey("error")) {
        throw Exception(responseJSON["error"]["message"]);
      }
    }
    throw Exception("Unrecognizable backend error");
  }

  Future<String> resetAPIKey(String name) async {
    var response = await client.post(
      Uri(
        scheme: kBackendScheme,
        host: kBackendHost,
        port: kBackendPort,
        pathSegments: [...relativePathSegments, name, "api-key"],
      ),
    );
    if (response.statusCode != 200) {
      handleErrorResponse(response);
    }
    return response.body;
  }
}
