const kBackendScheme =
    String.fromEnvironment("portal_backend_scheme", defaultValue: "http");

const kBackendHost =
    String.fromEnvironment("portal_backend_host", defaultValue: "localhost");

const kBackendPort =
    int.fromEnvironment("portal_backend_port", defaultValue: 8080);

const kBackendPath =
    String.fromEnvironment("portal_backend_path", defaultValue: "/internal");

final List<String> kBackendSegments = () {
  final matches = RegExp(r"(?:\/)?([^\/]+)").allMatches(kBackendPath);
  List<String> result = [];
  for (final m in matches) {
    result.add(m.group(1)!);
  }
  return List<String>.unmodifiable(result);
}();

const double maxDisplayWidth = 1000;
