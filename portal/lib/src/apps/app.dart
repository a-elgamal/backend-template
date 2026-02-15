import 'package:portal/src/stored/json.dart';

class App extends JSONSerializable {
  static const apiKeyJSONKey = "apiKey";
  static const disabledJSONKey = "disabled";

  final String apiKey;
  final bool disabled;

  const App({
    required this.apiKey,
    this.disabled = false,
  });

  App.fromJSON(Map<String, dynamic> json)
      : apiKey = json[apiKeyJSONKey],
        disabled = json[disabledJSONKey] ?? false;

  App.copy(App app, {String? apiKey, bool? disabled})
      : apiKey = apiKey ?? app.apiKey,
        disabled = disabled ?? app.disabled;

  @override
  Map<String, dynamic> toJSON() {
    return <String, dynamic>{
      App.apiKeyJSONKey: apiKey,
      App.disabledJSONKey: disabled,
    };
  }

  @override
  bool operator ==(Object other) =>
      other is App &&
      other.runtimeType == runtimeType &&
      other.apiKey == apiKey &&
      other.disabled == disabled;

  @override
  int get hashCode => Object.hash(apiKey, disabled);
}
