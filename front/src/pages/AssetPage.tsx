import { useState } from "react";
import { Link as RouterLink } from "react-router-dom";
import { useQuery, useMutation } from "urql";
import { graphql } from "../gql";
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
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  IconButton,
  List,
  ListItem,
  ListItemText,
  Divider,
  Card,
  CardContent,
  Alert
} from "@mui/material";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import AddIcon from "@mui/icons-material/Add";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";

// GraphQL クエリ定義
const GetAssetsDocument = graphql(/* GraphQL */ `
  query GetAssets {
    assets(first: 100) {
      nodes {
        id
        name
        category {
            id
            name
        }
      }
      pageInfo {
        hasNextPage
        endCursor   
        startCursor
        endCursor
      }
    }
  }
`);

// 資産作成のためのミューテーション
const CreateAssetDocument = graphql(/* GraphQL */ `
  mutation CreateAsset($name: String!, $categoryId: ID) {
    createAsset(input: { name: $name, categoryId: $categoryId }) {
      id
      name
      category {
        id
        name
      }
    }
  }
`);

// 資産更新のためのミューテーション
const UpdateAssetDocument = graphql(/* GraphQL */ `
  mutation UpdateAsset($id: ID!, $name: String!, $categoryId: ID) {
    updateAsset(input: { id: $id, name: $name, categoryId: $categoryId }) {
      id
      name
      category {
        id
        name
      }
    }
  }
`);

// 資産削除のためのミューテーション
const DeleteAssetDocument = graphql(/* GraphQL */ `
  mutation DeleteAsset($id: ID!) {
    deleteAsset(id: $id) {
      id
    }
  }
`);

export default function AssetPage() {
  // 資産フォームの状態
  const [assetForm, setAssetForm] = useState({
    id: "",
    name: "",
    categoryId: ""
  });

  // ダイアログの状態
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedAssetId, setSelectedAssetId] = useState("");
  
  // エラー通知の状態
  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");

  // GraphQL クエリの実行
  const [{ data, fetching, error }, reexecuteQuery] = useQuery({ query: GetAssetsDocument });
  
  // ミューテーション関数の取得
  const [, createAsset] = useMutation(CreateAssetDocument);
  const [, updateAsset] = useMutation(UpdateAssetDocument);
  const [, deleteAssetMutation] = useMutation(DeleteAssetDocument);

  // 作成ダイアログを開く
  const handleOpenCreateDialog = () => {
    setAssetForm({
      id: "",
      name: "",
      categoryId: ""
    });
    setCreateDialogOpen(true);
  };

  // 編集ダイアログを開く
  const handleOpenEditDialog = (asset: { id: string, name: string, category?: { id: string } | null }) => {
    setAssetForm({
      id: asset.id,
      name: asset.name,
      categoryId: asset.category?.id || ""
    });
    setSelectedAssetId(asset.id);
    setEditDialogOpen(true);
  };

  // 削除ダイアログを開く
  const handleOpenDeleteDialog = (assetId: string) => {
    setSelectedAssetId(assetId);
    setDeleteDialogOpen(true);
  };

  // フォームの変更を処理
  const handleFormChange = (field: string, value: string) => {
    setAssetForm(prev => ({ ...prev, [field]: value }));
  };

  // 資産を作成
  const handleCreateAsset = async () => {
    const variables = {
      name: assetForm.name,
      categoryId: assetForm.categoryId || null
    };
    
    try {
      const result = await createAsset(variables);
      if (result.error) {
        setErrorMessage(result.error.message);
      } else {
        setSuccessMessage("資産が正常に作成されました");
        setCreateDialogOpen(false);
        reexecuteQuery({ requestPolicy: 'network-only' });
      }
    } catch (err) {
      setErrorMessage("資産の作成中にエラーが発生しました");
    }
  };

  // 資産を更新
  const handleUpdateAsset = async () => {
    const variables = {
      id: assetForm.id,
      name: assetForm.name,
      categoryId: assetForm.categoryId || null
    };
    
    try {
      const result = await updateAsset(variables);
      if (result.error) {
        setErrorMessage(result.error.message);
      } else {
        setSuccessMessage("資産が正常に更新されました");
        setEditDialogOpen(false);
        reexecuteQuery({ requestPolicy: 'network-only' });
      }
    } catch (err) {
      setErrorMessage("資産の更新中にエラーが発生しました");
    }
  };

  // 資産を削除
  const handleDeleteAsset = async () => {
    try {
      const result = await deleteAssetMutation({ id: selectedAssetId });
      if (result.error) {
        setErrorMessage(result.error.message);
      } else {
        setSuccessMessage("資産が正常に削除されました");
        setDeleteDialogOpen(false);
        reexecuteQuery({ requestPolicy: 'network-only' });
      }
    } catch (err) {
      setErrorMessage("資産の削除中にエラーが発生しました");
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

  const assets = data?.assets?.nodes || [];
  // サーバーからカテゴリデータを取得するクエリがないため、資産データから取得
  const assetCategoriesMap = new Map();
  
  // 資産からカテゴリ情報を抽出
  assets.forEach(asset => {
    if (asset.category) {
      assetCategoriesMap.set(asset.category.id, asset.category);
    }
  });
  
  const assetCategories = Array.from(assetCategoriesMap.values());

  // カテゴリでグループ化された資産を取得
  interface CategoryGroup {
    id: string;
    name: string;
    assets: Array<{
      id: string;
      name: string;
      category?: { id: string; name: string } | null;
    }>;
  }
  
  const assetsByCategory = assets.reduce<Record<string, CategoryGroup>>((acc, asset) => {
    const categoryId = asset.category?.id || "uncategorized";
    const categoryName = asset.category?.name || "未分類";
    
    if (!acc[categoryId]) {
      acc[categoryId] = {
        id: categoryId,
        name: categoryName,
        assets: []
      };
    }
    
    acc[categoryId].assets.push(asset);
    return acc;
  }, {});

  return (
    <Box p={3}>
      {/* ヘッダー */}
      <Stack direction="row" alignItems="center" justifyContent="space-between" mb={3}>
        <Button component={RouterLink} to="/" startIcon={<ArrowBackIcon />}>
          ホームに戻る
        </Button>
        <Typography variant="h5">資産一覧</Typography>
        <Button 
          variant="contained" 
          color="primary" 
          startIcon={<AddIcon />}
          onClick={handleOpenCreateDialog}
        >
          新規資産
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

      {/* 資産カテゴリへのリンク */}
      <Box mb={3}>
        <Button 
          component={RouterLink} 
          to="/asset/category" 
          variant="outlined" 
          color="primary"
        >
          資産カテゴリ管理
        </Button>
      </Box>

      {/* 資産リスト */}
      {Object.values(assetsByCategory).length === 0 ? (
        <Paper elevation={1} sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="h6" color="text.secondary">
            資産がありません
          </Typography>
          <Button 
            variant="contained" 
            color="primary"
            startIcon={<AddIcon />} 
            onClick={handleOpenCreateDialog}
            sx={{ mt: 2 }}
          >
            新規資産を作成
          </Button>
        </Paper>
      ) : (
        <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(12, 1fr)', gap: 3 }}>
          {Object.values(assetsByCategory).map(category => (
            <Box sx={{ gridColumn: {xs: 'span 12', md: 'span 6'} }} key={category.id}>
              <Card variant="outlined">
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    {category.name}
                  </Typography>
                  <Divider sx={{ mb: 2 }} />
                  <List dense disablePadding>
                    {category.assets.map((asset) => (
                      <ListItem
                        key={asset.id}
                        secondaryAction={
                          <Stack direction="row" spacing={1}>
                            <IconButton edge="end" onClick={() => handleOpenEditDialog(asset)}>
                              <EditIcon fontSize="small" />
                            </IconButton>
                            <IconButton edge="end" onClick={() => handleOpenDeleteDialog(asset.id)}>
                              <DeleteIcon fontSize="small" />
                            </IconButton>
                          </Stack>
                        }
                      >
                        <ListItemText primary={asset.name} />
                      </ListItem>
                    ))}
                  </List>
                </CardContent>
              </Card>
            </Box>
          ))}
        </Box>
      )}

      {/* 新規資産作成ダイアログ */}
      <Dialog open={createDialogOpen} onClose={() => setCreateDialogOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>新規資産作成</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="資産名"
            fullWidth
            variant="outlined"
            value={assetForm.name}
            onChange={(e) => handleFormChange('name', e.target.value)}
            sx={{ mb: 2, mt: 2 }}
          />
          <FormControl fullWidth>
            <InputLabel>カテゴリ</InputLabel>
            <Select
              value={assetForm.categoryId}
              onChange={(e) => handleFormChange('categoryId', e.target.value)}
              label="カテゴリ"
            >
              <MenuItem value="">未分類</MenuItem>
              {assetCategories.map((category) => (
                <MenuItem key={category.id} value={category.id}>
                  {category.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateDialogOpen(false)}>キャンセル</Button>
          <Button 
            onClick={handleCreateAsset} 
            variant="contained" 
            color="primary"
            disabled={!assetForm.name.trim()}
          >
            作成
          </Button>
        </DialogActions>
      </Dialog>

      {/* 資産編集ダイアログ */}
      <Dialog open={editDialogOpen} onClose={() => setEditDialogOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>資産を編集</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="資産名"
            fullWidth
            variant="outlined"
            value={assetForm.name}
            onChange={(e) => handleFormChange('name', e.target.value)}
            sx={{ mb: 2, mt: 2 }}
          />
          <FormControl fullWidth>
            <InputLabel>カテゴリ</InputLabel>
            <Select
              value={assetForm.categoryId}
              onChange={(e) => handleFormChange('categoryId', e.target.value)}
              label="カテゴリ"
            >
              <MenuItem value="">未分類</MenuItem>
              {assetCategories.map((category) => (
                <MenuItem key={category.id} value={category.id}>
                  {category.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialogOpen(false)}>キャンセル</Button>
          <Button 
            onClick={handleUpdateAsset} 
            variant="contained" 
            color="primary"
            disabled={!assetForm.name.trim()}
          >
            更新
          </Button>
        </DialogActions>
      </Dialog>

      {/* 資産削除確認ダイアログ */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>資産を削除</DialogTitle>
        <DialogContent>
          <Typography>本当にこの資産を削除しますか？この操作は元に戻せません。</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>キャンセル</Button>
          <Button onClick={handleDeleteAsset} variant="contained" color="error">
            削除
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
