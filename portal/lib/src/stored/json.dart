abstract class JSONSerializable {
  const JSONSerializable();

  Map<String, dynamic> toJSON();

  @override
  String toString() {
    return toJSON().toString();
  }
}
