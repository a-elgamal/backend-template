import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class NewAppView extends StatefulWidget {
  final Function(String) onSave;

  const NewAppView({
    super.key,
    required this.onSave,
  });

  @override
  NewAppViewState createState() => NewAppViewState();
}

class NewAppViewState extends State<NewAppView> {
  final _textController = TextEditingController();
  static const maxNameLength = 36;

  @override
  void initState() {
    super.initState();

    _textController.addListener(() {
      setState(() {});
    });
  }

  @override
  void dispose() {
    _textController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        TextField(
          controller: _textController,
          autofocus: true,
          maxLength: maxNameLength,
          inputFormatters: [
            FilteringTextInputFormatter.allow(RegExp("[0-9a-z]")),
          ],
          decoration: InputDecoration(
            label: const Text("App Name"),
            counterText: "${_textController.text.length}/$maxNameLength",
          ),
          onSubmitted: (name) => widget.onSave(name),
        ),
        const SizedBox(
          height: 5,
        ),
        Row(
          children: [
            const Spacer(),
            ElevatedButton(
              onPressed: _textController.text.isNotEmpty
                  ? () => widget.onSave(_textController.text)
                  : null,
              child: const Text("Save"),
            ),
          ],
        )
      ],
    );
  }
}
