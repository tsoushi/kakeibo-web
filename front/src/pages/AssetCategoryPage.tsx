import { useState } from "react";
import { Link as RouterLink } from "react-router-dom";
import { useQuery, useMutation } from "urql";
import { graphql } from "../gql";
import type { AssetCategory, GetAssetCategoriesQuery } from "../gql/graphql";
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
  Alert
} from "@mui/material";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import AddIcon from "@mui/icons-material/Add";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";

// GraphQL クエリ定義
const GetAssetCategoriesDocument = graphql(/* GraphQL */ `
  query GetAssetCategories {
    assetCategories(first: 100) {
        nodes {
            id
            name
            assets {
                id
                name
            }
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

// 資産カテゴリ作成のためのミューテーション
const CreateAssetCategoryDocument = graphql(/* GraphQL */ `
  mutation CreateAssetCategory($name: String!) {
    createAssetCategory(input: { name: $name }) {
      id
      name
    }
  }
`);

// 資産カテゴリ更新のためのミューテーション
const UpdateAssetCategoryDocument = graphql(/* GraphQL */ `
  mutation UpdateAssetCategory($id: ID!, $name: String!) {
    updateAssetCategory(input: { id: $id, name: $name }) {
      id
      name
    }
  }
`);

// 資産カテゴリ削除のためのミューテーション
const DeleteAssetCategoryDocument = graphql(/* GraphQL */ `
  mutation DeleteAssetCategory($id: ID!) {
    deleteAssetCategory(input: { id: $id }) {
        id
    }
  }
`);

export default function AssetCategoryPage() {
  // カテゴリフォームの状態
  const [categoryForm, setCategoryForm] = useState({
    id: "",
    name: ""
  });

  // ダイアログの状態
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedCategoryId, setSelectedCategoryId] = useState("");
  const [selectedCategoryName, setSelectedCategoryName] = useState("");
  const [selectedCategoryHasAssets, setSelectedCategoryHasAssets] = useState(false);
  
  // エラー通知の状態
  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");

  // GraphQL クエリの実行
  const [{ data, fetching, error }, reexecuteQuery] = useQuery({ query: GetAssetCategoriesDocument });
  
  // ミューテーション関数の取得
  const [, createCategory] = useMutation(CreateAssetCategoryDocument);
  const [, updateCategory] = useMutation(UpdateAssetCategoryDocument);
  const [, deleteCategory] = useMutation(DeleteAssetCategoryDocument);

  // 作成ダイアログを開く
  const handleOpenCreateDialog = () => {
    setCategoryForm({
      id: "",
      name: ""
    });
    setCreateDialogOpen(true);
  };

  // 編集ダイアログを開く
  const handleOpenEditDialog = (category: AssetCategory) => {
    setCategoryForm({
      id: category.id,
      name: category.name
    });
    setSelectedCategoryId(category.id);
    setEditDialogOpen(true);
  };

  // 削除ダイアログを開く
  const handleOpenDeleteDialog = (category: GetAssetCategoriesQuery["assetCategories"]["nodes"][number]) => {
    setSelectedCategoryId(category.id);
    setSelectedCategoryName(category.name);
    setSelectedCategoryHasAssets(category.assets && category.assets.length > 0);
    setDeleteDialogOpen(true);
  };

  // フォームの変更を処理
  const handleFormChange = (field: string, value: string) => {
    setCategoryForm(prev => ({ ...prev, [field]: value }));
  };

  // カテゴリを作成
  const handleCreateCategory = async () => {
    const variables = {
      name: categoryForm.name
    };
    
    try {
      const result = await createCategory(variables);
      if (result.error) {
        setErrorMessage(result.error.message);
      } else {
        setSuccessMessage("カテゴリが正常に作成されました");
        setCreateDialogOpen(false);
        reexecuteQuery({ requestPolicy: 'network-only' });
      }
    } catch (err) {
      setErrorMessage("カテゴリの作成中にエラーが発生しました");
    }
  };

  // カテゴリを更新
  const handleUpdateCategory = async () => {
    const variables = {
      id: categoryForm.id,
      name: categoryForm.name
    };
    
    try {
      const result = await updateCategory(variables);
      if (result.error) {
        setErrorMessage(result.error.message);
      } else {
        setSuccessMessage("カテゴリが正常に更新されました");
        setEditDialogOpen(false);
        reexecuteQuery({ requestPolicy: 'network-only' });
      }
    } catch (err) {
      setErrorMessage("カテゴリの更新中にエラーが発生しました");
    }
  };

  // カテゴリを削除
  const handleDeleteCategory = async () => {
    try {
      const result = await deleteCategory({ id: selectedCategoryId });
      if (result.error) {
        setErrorMessage(result.error.message);
      } else {
        setSuccessMessage("カテゴリが正常に削除されました");
        setDeleteDialogOpen(false);
        reexecuteQuery({ requestPolicy: 'network-only' });
      }
    } catch (err) {
      setErrorMessage("カテゴリの削除中にエラーが発生しました");
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

  const categories = data?.assetCategories.nodes || [];

  return (
    <Box p={3}>
      {/* ヘッダー */}
      <Stack direction="row" alignItems="center" justifyContent="space-between" mb={3}>
        <Button component={RouterLink} to="/asset" startIcon={<ArrowBackIcon />}>
          資産一覧に戻る
        </Button>
        <Typography variant="h5">資産カテゴリ一覧</Typography>
        <Button 
          variant="contained" 
          color="primary" 
          startIcon={<AddIcon />}
          onClick={handleOpenCreateDialog}
        >
          新規カテゴリ
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

      {/* カテゴリリスト */}
      {categories.length === 0 ? (
        <Paper elevation={1} sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="h6" color="text.secondary">
            資産カテゴリがありません
          </Typography>
          <Button 
            variant="contained" 
            color="primary"
            startIcon={<AddIcon />} 
            onClick={handleOpenCreateDialog}
            sx={{ mt: 2 }}
          >
            新規カテゴリを作成
          </Button>
        </Paper>
      ) : (
        <Paper elevation={1} sx={{ mt: 2 }}>
          <List>
            {categories.map((category, index) => (
              <div key={category.id}>
                {index > 0 && <Divider />}
                <ListItem
                  secondaryAction={
                    <Stack direction="row" spacing={1}>
                      <IconButton edge="end" onClick={() => handleOpenEditDialog(category)}>
                        <EditIcon />
                      </IconButton>
                      <IconButton edge="end" onClick={() => handleOpenDeleteDialog(category)}>
                        <DeleteIcon />
                      </IconButton>
                    </Stack>
                  }
                >
                  <ListItemText 
                    primary={category.name} 
                    secondary={
                      category.assets && category.assets.length > 0 
                        ? `関連する資産: ${category.assets.length}個` 
                        : "関連する資産なし"
                    } 
                  />
                </ListItem>
              </div>
            ))}
          </List>
        </Paper>
      )}

      {/* 新規カテゴリ作成ダイアログ */}
      <Dialog open={createDialogOpen} onClose={() => setCreateDialogOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>新規カテゴリ作成</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="カテゴリ名"
            fullWidth
            variant="outlined"
            value={categoryForm.name}
            onChange={(e) => handleFormChange('name', e.target.value)}
            sx={{ mb: 1, mt: 2 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateDialogOpen(false)}>キャンセル</Button>
          <Button 
            onClick={handleCreateCategory} 
            variant="contained" 
            color="primary"
            disabled={!categoryForm.name.trim()}
          >
            作成
          </Button>
        </DialogActions>
      </Dialog>

      {/* カテゴリ編集ダイアログ */}
      <Dialog open={editDialogOpen} onClose={() => setEditDialogOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>カテゴリを編集</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="カテゴリ名"
            fullWidth
            variant="outlined"
            value={categoryForm.name}
            onChange={(e) => handleFormChange('name', e.target.value)}
            sx={{ mb: 1, mt: 2 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialogOpen(false)}>キャンセル</Button>
          <Button 
            onClick={handleUpdateCategory} 
            variant="contained" 
            color="primary"
            disabled={!categoryForm.name.trim()}
          >
            更新
          </Button>
        </DialogActions>
      </Dialog>

      {/* カテゴリ削除確認ダイアログ */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>カテゴリを削除</DialogTitle>
        <DialogContent>
          {selectedCategoryHasAssets ? (
            <Typography color="error">
              このカテゴリ「{selectedCategoryName}」には関連する資産があります。
              先に資産のカテゴリを変更するか、資産を削除してから再試行してください。
            </Typography>
          ) : (
            <Typography>
              カテゴリ「{selectedCategoryName}」を本当に削除しますか？
              この操作は元に戻せません。
            </Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>キャンセル</Button>
          <Button 
            onClick={handleDeleteCategory} 
            variant="contained" 
            color="error"
            disabled={selectedCategoryHasAssets}
          >
            削除
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
