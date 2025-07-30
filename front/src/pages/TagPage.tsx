import { useState } from "react";
import { Link as RouterLink } from "react-router-dom";
import { useQuery, useMutation } from "urql";
import { graphql } from "../gql";
import type { Tag } from "../gql/graphql";
import {
  Box,
  Typography,
  Paper,
  CircularProgress,
  Button,
  Stack,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  List,
  ListItem,
  ListItemText,
  IconButton,
  Divider,
  Alert,
  Chip
} from "@mui/material";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import AddIcon from "@mui/icons-material/Add";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";

// GraphQL クエリ定義
const GetTagsDocument = graphql(/* GraphQL */ `
  query GetTags {
    tags(first: 100) {
      nodes {
        id
        name
      }
      pageInfo {
        hasNextPage
        hasPreviousPage
        startCursor
        endCursor
      }
    }
  }
`);

// タグ作成のためのミューテーション
const CreateTagDocument = graphql(/* GraphQL */ `
  mutation CreateTag($name: String!) {
    createTag(input: { name: $name }) {
      id
      name
    }
  }
`);

// タグ更新のためのミューテーション
const UpdateTagDocument = graphql(/* GraphQL */ `
  mutation UpdateTag($id: ID!, $name: String!) {
    updateTag(input: { id: $id, name: $name }) {
      id
      name
    }
  }
`);

// タグ削除のためのミューテーション
const DeleteTagDocument = graphql(/* GraphQL */ `
  mutation DeleteTag($id: ID!) {
    deleteTag(input: { id: $id }) {
      id
    }
  }
`);

export default function TagPage() {
  // タグフォームの状態
  const [tagForm, setTagForm] = useState({
    id: "",
    name: ""
  });

  // ダイアログの状態
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedTagId, setSelectedTagId] = useState("");
  const [selectedTagName, setSelectedTagName] = useState("");
  
  // エラー通知の状態
  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");

  // GraphQL クエリの実行
  const [{ data, fetching, error }, reexecuteQuery] = useQuery({ query: GetTagsDocument });
  
  // ミューテーション関数の取得
  const [, createTag] = useMutation(CreateTagDocument);
  const [, updateTag] = useMutation(UpdateTagDocument);
  const [, deleteTagMutation] = useMutation(DeleteTagDocument);

  // 作成ダイアログを開く
  const handleOpenCreateDialog = () => {
    setTagForm({
      id: "",
      name: ""
    });
    setCreateDialogOpen(true);
  };

  // 編集ダイアログを開く
  const handleOpenEditDialog = (tag: Tag) => {
    setTagForm({
      id: tag.id,
      name: tag.name
    });
    setSelectedTagId(tag.id);
    setEditDialogOpen(true);
  };

  // 削除ダイアログを開く
  const handleOpenDeleteDialog = (tag: Tag) => {
    setSelectedTagId(tag.id);
    setSelectedTagName(tag.name);
    setDeleteDialogOpen(true);
  };

  // フォームの変更を処理
  const handleFormChange = (field: string, value: string) => {
    setTagForm(prev => ({ ...prev, [field]: value }));
  };

  // タグを作成
  const handleCreateTag = async () => {
    const variables = {
      name: tagForm.name
    };
    
    try {
      const result = await createTag(variables);
      if (result.error) {
        setErrorMessage(result.error.message);
      } else {
        setSuccessMessage("タグが正常に作成されました");
        setCreateDialogOpen(false);
        reexecuteQuery({ requestPolicy: 'network-only' });
      }
    } catch (err) {
      setErrorMessage("タグの作成中にエラーが発生しました");
    }
  };

  // タグを更新
  const handleUpdateTag = async () => {
    const variables = {
      id: tagForm.id,
      name: tagForm.name
    };
    
    try {
      const result = await updateTag(variables);
      if (result.error) {
        setErrorMessage(result.error.message);
      } else {
        setSuccessMessage("タグが正常に更新されました");
        setEditDialogOpen(false);
        reexecuteQuery({ requestPolicy: 'network-only' });
      }
    } catch (err) {
      setErrorMessage("タグの更新中にエラーが発生しました");
    }
  };

  // タグを削除
  const handleDeleteTag = async () => {
    try {
      const result = await deleteTagMutation({ id: selectedTagId });
      if (result.error) {
        setErrorMessage(result.error.message);
      } else {
        setSuccessMessage("タグが正常に削除されました");
        setDeleteDialogOpen(false);
        reexecuteQuery({ requestPolicy: 'network-only' });
      }
    } catch (err) {
      setErrorMessage("タグの削除中にエラーが発生しました");
    }
  };

  // メッセージをクリア
  const clearMessages = () => {
    setErrorMessage("");
    setSuccessMessage("");
  };

  if (fetching) {
    return (
      <Box p={6} textAlign="center">
        <CircularProgress />
        <Typography mt={2}>読み込み中...</Typography>
      </Box>
    );
  }

  if (error) {
    return (
      <Paper elevation={3} sx={{ p: 3, maxWidth: 600, mx: "auto", textAlign: "center" }}>
        <Typography variant="h6" color="error">エラー</Typography>
        <Typography mt={2}>{error.message}</Typography>
        <Button component={RouterLink} to="/" variant="contained" color="primary" sx={{ mt: 2 }}>
          ホームに戻る
        </Button>
      </Paper>
    );
  }

  const tags = data?.tags.nodes || [];

  return (
    <Box p={3}>
      {/* ヘッダー */}
      <Stack direction="row" alignItems="center" justifyContent="space-between" mb={3}>
        <Button component={RouterLink} to="/" startIcon={<ArrowBackIcon />}>
          ホームに戻る
        </Button>
        <Typography variant="h5">タグ一覧</Typography>
        <Button 
          variant="contained" 
          color="primary" 
          startIcon={<AddIcon />}
          onClick={handleOpenCreateDialog}
        >
          新規タグ
        </Button>
      </Stack>

      {/* メッセージ表示 */}
      {errorMessage && (
        <Alert severity="error" onClose={clearMessages} sx={{ mb: 2 }}>
          {errorMessage}
        </Alert>
      )}
      {successMessage && (
        <Alert severity="success" onClose={clearMessages} sx={{ mb: 2 }}>
          {successMessage}
        </Alert>
      )}

      {/* タグリスト */}
      {tags.length === 0 ? (
        <Paper elevation={1} sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="h6" color="text.secondary">
            タグがありません
          </Typography>
          <Button 
            variant="contained" 
            color="primary"
            startIcon={<AddIcon />} 
            onClick={handleOpenCreateDialog}
            sx={{ mt: 2 }}
          >
            新規タグを作成
          </Button>
        </Paper>
      ) : (
        <Box sx={{ mt: 2 }}>
          <Typography variant="subtitle1" gutterBottom>すべてのタグ</Typography>
          <Paper elevation={1} sx={{ p: 3 }}>
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
              {tags.map((tag) => (
                <Chip 
                  key={tag.id}
                  label={tag.name}
                  color="primary"
                  onDelete={() => handleOpenDeleteDialog(tag)}
                  onClick={() => handleOpenEditDialog(tag)}
                  sx={{ m: 0.5 }}
                />
              ))}
            </Box>
          </Paper>
          
          <Paper elevation={1} sx={{ mt: 3 }}>
            <List>
              {tags.map((tag, index) => (
                <div key={tag.id}>
                  {index > 0 && <Divider />}
                  <ListItem
                    secondaryAction={
                      <Stack direction="row" spacing={1}>
                        <IconButton edge="end" onClick={() => handleOpenEditDialog(tag)}>
                          <EditIcon />
                        </IconButton>
                        <IconButton edge="end" onClick={() => handleOpenDeleteDialog(tag)}>
                          <DeleteIcon />
                        </IconButton>
                      </Stack>
                    }
                  >
                    <ListItemText primary={tag.name} />
                  </ListItem>
                </div>
              ))}
            </List>
          </Paper>
        </Box>
      )}

      {/* 新規タグ作成ダイアログ */}
      <Dialog open={createDialogOpen} onClose={() => setCreateDialogOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>新規タグ作成</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="タグ名"
            fullWidth
            variant="outlined"
            value={tagForm.name}
            onChange={(e) => handleFormChange('name', e.target.value)}
            sx={{ mb: 1, mt: 2 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateDialogOpen(false)}>キャンセル</Button>
          <Button 
            onClick={handleCreateTag} 
            variant="contained" 
            color="primary"
            disabled={!tagForm.name.trim()}
          >
            作成
          </Button>
        </DialogActions>
      </Dialog>

      {/* タグ編集ダイアログ */}
      <Dialog open={editDialogOpen} onClose={() => setEditDialogOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>タグを編集</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="タグ名"
            fullWidth
            variant="outlined"
            value={tagForm.name}
            onChange={(e) => handleFormChange('name', e.target.value)}
            sx={{ mb: 1, mt: 2 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialogOpen(false)}>キャンセル</Button>
          <Button 
            onClick={handleUpdateTag} 
            variant="contained" 
            color="primary"
            disabled={!tagForm.name.trim()}
          >
            更新
          </Button>
        </DialogActions>
      </Dialog>

      {/* タグ削除確認ダイアログ */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>タグを削除</DialogTitle>
        <DialogContent>
          <Typography>
            タグ「{selectedTagName}」を本当に削除しますか？
            この操作は元に戻せません。
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>キャンセル</Button>
          <Button 
            onClick={handleDeleteTag} 
            variant="contained" 
            color="error"
          >
            削除
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
