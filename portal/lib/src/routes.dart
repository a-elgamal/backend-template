import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:portal/src/apps/app.dart';
import 'package:portal/src/apps/app_list_view.dart';
import 'package:portal/src/apps/app_service.dart';
import 'package:portal/src/apps/app_view.dart';
import 'package:portal/src/stored/stored.dart';
import 'package:provider/provider.dart';

part 'routes.g.dart';

// After changing this file, run `flutter pub run build_running build` to generate routes.g.dart.

@TypedGoRoute<AppListViewRoute>(
  name: 'appListView',
  path: '/apps',
  routes: <TypedGoRoute<GoRouteData>>[
    TypedGoRoute<AppViewRoute>(
      name: 'appView',
      path: ":appName",
    )
  ],
)
class AppListViewRoute extends GoRouteData with $AppListViewRoute {
  const AppListViewRoute();

  @override
  Widget build(BuildContext context, GoRouterState state) =>
      Consumer<AppService>(
        builder: (_, s, _) => AppListView(s),
      );
}

class AppViewRoute extends GoRouteData with $AppViewRoute {
  final String appName;
  final Stored<App>? $extra;

  const AppViewRoute(this.appName, this.$extra);

  @override
  Widget build(BuildContext context, GoRouterState state) =>
      Consumer<AppService>(
        builder: (_, s, _) => AppView(
          s,
          appName,
          app: $extra,
        ),
      );
}
