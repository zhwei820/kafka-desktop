<template>
  <!-- theme="vs-dark" -->

  <MonacoEditor
    :code="code"
    height="300"
    theme="vs-light"
    :editor-options="MONACO_EDITOR_OPTIONS"
    @mounted="onMounted"
    @codeChange="onCodeChange"
    @selectionChange="onSelectionChange"
    language="sql"
  />
</template>

<script>
import MonacoEditor from "../comp/vue-monaco-editor";
import { EventsEmit } from "../../wailsjs/runtime";
import { format } from "sql-formatter";

let defaultSuggestions = [];

export default {
  components: {
    MonacoEditor,
  },
  props: {
    keywords: {
      type: Array, // Specify the type (e.g., Array)
      required: true, // Optional, makes the prop required
    },
    value: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      MONACO_EDITOR_OPTIONS: {
        automaticLayout: true,
        formatOnType: true,
        formatOnPaste: true,
      },
      code: this.value,
      editor: {},
    };
  },
  computed: {},
  watch: {
    keywords() {
      if (window.monaco) {
        this.updateKeywords();
      }
    },
  },

  methods: {
    onMounted(editor) {
      this.editor = editor;

      let suggestions = function (range) {
        defaultSuggestions = [
          {
            label: "select *",
            kind: window.monaco.languages.CompletionItemKind.Property,
            range: range,
            insertText: "select *",
          },
          {
            label: "update",
            kind: window.monaco.languages.CompletionItemKind.Property,
            insertText: "update ",
          },
          {
            label: "delete",
            kind: window.monaco.languages.CompletionItemKind.Property,
            insertText: "delete ",
          },
          {
            label: "alter",
            kind: window.monaco.languages.CompletionItemKind.Property,
            insertText: "alter ",
          },
          {
            label: "from",
            kind: window.monaco.languages.CompletionItemKind.Property,
            insertText: "from ",
          },
          {
            label: "where",
            kind: window.monaco.languages.CompletionItemKind.Property,
            insertText: "where ",
          },
          {
            label: "order by",
            kind: window.monaco.languages.CompletionItemKind.Property,
            insertText: "order by ",
          },
          {
            label: "limit",
            kind: window.monaco.languages.CompletionItemKind.Property,
            insertText: "limit ",
          },
          {
            label: "now()",
            kind: window.monaco.languages.CompletionItemKind.Property,
            insertText: "now() ",
          },
          {
            label: "and",
            kind: window.monaco.languages.CompletionItemKind.Property,
            insertText: "and ",
          },
          {
            label: "or",
            kind: window.monaco.languages.CompletionItemKind.Property,
            insertText: "or ",
          },
          {
            label: "id",
            kind: window.monaco.languages.CompletionItemKind.Property,
            insertText: "id ",
          },
          {
            label: "from_unixtime()",
            kind: window.monaco.languages.CompletionItemKind.Property,
            insertText: "from_unixtime() ",
          },
          {
            label: "unix_timestamp()",
            kind: window.monaco.languages.CompletionItemKind.Property,
            insertText: "unix_timestamp() ",
          },
        ];
        return [...defaultSuggestions];
      };
      // console.log("==>suggestions", suggestions);

      window.monaco.languages.registerCompletionItemProvider("sql", {
        provideCompletionItems: function (model, position) {
          var word = model.getWordUntilPosition(position);
          var range = {
            startLineNumber: position.lineNumber,
            endLineNumber: position.lineNumber,
            startColumn: word.startColumn,
            endColumn: word.endColumn,
          };

          return {
            suggestions: suggestions(range),
          };
        },
      });

      this.updateKeywords();

      this.editor.addAction({
        // An unique identifier of the contributed action.
        id: "my-unique-id",

        // A label of the action that will be presented to the user.
        label: "Run",

        // An optional array of keybindings for the action.
        keybindings: [
          window.monaco.KeyMod.CtrlCmd | window.monaco.KeyCode.Enter,
        ],

        // A precondition for this action.
        precondition: null,

        // A rule to evaluate on top of the precondition in order to dispatch the keybindings.
        keybindingContext: null,

        contextMenuGroupId: "navigation",

        contextMenuOrder: 1.5,

        // Method that will be executed when the action is triggered.
        // @param editor The editor instance is passed in as a convenience
        run: function (ed) {
          EventsEmit("onQuery");
        },
      });
      let that = this;
      this.editor.addAction({
        // An unique identifier of the contributed action.
        id: "my-unique-id2",

        // A label of the action that will be presented to the user.
        label: "Format",

        // An optional array of keybindings for the action.
        keybindings: [
          // window.monaco.KeyMod.CtrlCmd | window.monaco.KeyCode.KeyF,
        ],

        // A precondition for this action.
        precondition: null,

        // A rule to evaluate on top of the precondition in order to dispatch the keybindings.
        keybindingContext: null,

        contextMenuGroupId: "navigation",

        contextMenuOrder: 1.5,

        // Method that will be executed when the action is triggered.
        // @param editor The editor instance is passed in as a convenience
        run: function (ed) {
          let code = that.editor.getValue();
          code = format(code, { linesBetweenQueries: 2 });
          that.editor.setValue(code);
        },
      });
    },
    updateKeywords() {
      let keywords = this.keywords;
      let suggestions = function (range) {
        return [
          ...keywords.map((val) => {
            return {
              label: val,
              kind: window.monaco.languages.CompletionItemKind.Property,
              insertText: val,
            };
          }),
        ];
      };
      // console.log("==>suggestions", suggestions);

      window.monaco.languages.registerCompletionItemProvider("sql", {
        provideCompletionItems: function (model, position) {
          var word = model.getWordUntilPosition(position);
          var range = {
            startLineNumber: position.lineNumber,
            endLineNumber: position.lineNumber,
            startColumn: word.startColumn,
            endColumn: word.endColumn,
          };

          return {
            suggestions: suggestions(range),
          };
        },
      });
    },
    onCodeChange(editor) {
      let code = editor.getValue();
      // console.log("===>code", code);
      this.$emit("onUpdateValue", code);
    },
    onSelectionChange(editor) {
      //   let editor = this.editor;
      var selectedText = editor
        .getModel()
        .getValueInRange(editor.getSelection());

      //   console.log("selectedText", selectedText);
      this.$emit("onSelectionChange", selectedText);
      return selectedText;
    },
    formatCode: function () {
      this.editor?.getAction("editor.action.formatDocument").run();
    },
  },

  mounted() {},
};

// your action
</script>
