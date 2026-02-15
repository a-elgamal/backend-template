// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'routes.dart';

// **************************************************************************
// GoRouterGenerator
// **************************************************************************

List<RouteBase> get $appRoutes => [$appListViewRoute];

RouteBase get $appListViewRoute => GoRouteData.$route(
  path: '/apps',
  name: 'appListView',
  factory: $AppListViewRoute._fromState,
  routes: [
    GoRouteData.$route(
      path: ':appName',
      name: 'appView',
      factory: $AppViewRoute._fromState,
    ),
  ],
);

mixin $AppListViewRoute on GoRouteData {
  static AppListViewRoute _fromState(GoRouterState state) =>
      const AppListViewRoute();

  @override
  String get location => GoRouteData.$location('/apps');

  @override
  void go(BuildContext context) => context.go(location);

  @override
  Future<T?> push<T>(BuildContext context) => context.push<T>(location);

  @override
  void pushReplacement(BuildContext context) =>
      context.pushReplacement(location);

  @override
  void replace(BuildContext context) => context.replace(location);
}

mixin $AppViewRoute on GoRouteData {
  static AppViewRoute _fromState(GoRouterState state) => AppViewRoute(
    state.pathParameters['appName']!,
    state.extra as Stored<App>?,
  );

  AppViewRoute get _self => this as AppViewRoute;

  @override
  String get location =>
      GoRouteData.$location('/apps/${Uri.encodeComponent(_self.appName)}');

  @override
  void go(BuildContext context) => context.go(location, extra: _self.$extra);

  @override
  Future<T?> push<T>(BuildContext context) =>
      context.push<T>(location, extra: _self.$extra);

  @override
  void pushReplacement(BuildContext context) =>
      context.pushReplacement(location, extra: _self.$extra);

  @override
  void replace(BuildContext context) =>
      context.replace(location, extra: _self.$extra);
}
