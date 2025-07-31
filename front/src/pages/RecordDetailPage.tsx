import { useState } from "react";
import { useParams, Link as RouterLink, useNavigate } from "react-router-dom";
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
  Chip,
  Divider,
  Alert
} from "@mui/material";
import { DateTimePicker } from "@mui/x-date-pickers/DateTimePicker";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFnsV3";
import { ja } from "date-fns/locale";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";

// GraphQL クエリ定義
const GetRecordDocument = graphql(/* GraphQL */ `
  query GetRecord($id: ID!) {
    record(id: $id) {
      id
      recordType
      title
      description
      at
      assetChangeIncome {
        asset {
          id
          name
        }
        amount
      }
      assetChangeExpense {
        asset {
          id
          name
        }
        amount
      }
      tags {
        id
        name
      }
    }
  }
`);

// レコード更新のためのミューテーション
const UpdateIncomeRecordDocument = graphql(/* GraphQL */ `
  mutation UpdateIncomeRecord(
    $id: ID!,
    $title: String!,
    $description: String!,
    $at: Time!,
    $assetID: ID!,
    $amount: Int!,
    $tags: [String!]!
  ) {
    updateIncomeRecord(
      input: {
        id: $id,
        title: $title,
        description: $description,
        at: $at,
        assetID: $assetID,
        amount: $amount,
        tags: $tags
      }
    ) {
      id
      recordType
      title
      description
      at
      assetChangeIncome {
        asset {
          id
          name
        }
        amount
      }
    }
  }
`);

const UpdateExpenseRecordDocument = graphql(/* GraphQL */ `
  mutation UpdateExpenseRecord(
    $id: ID!,
    $title: String!,
    $description: String!,
    $at: Time!,
    $assetID: ID!,
    $amount: Int!,
    $tags: [String!]!
  ) {
    updateExpenseRecord(
      input: {
        id: $id,
        title: $title,
        description: $description,
        at: $at,
        assetID: $assetID,
        amount: $amount,
        tags: $tags
      }
    ) {
      id
      recordType
      title
      description
      at
      assetChangeExpense {
        asset {
          id
          name
        }
        amount
      }
    }
  }
`);

const UpdateTransferRecordDocument = graphql(/* GraphQL */ `
  mutation UpdateTransferRecord(
    $id: ID!,
    $title: String!,
    $description: String!,
    $at: Time!,
    $fromAssetID: ID!,
    $toAssetID: ID!,
    $amount: Int!,
    $tags: [String!]!
  ) {
    updateTransferRecord(
      input: {
        id: $id,
        title: $title,
        description: $description,
        at: $at,
        fromAssetID: $fromAssetID,
        toAssetID: $toAssetID,
        amount: $amount,
        tags: $tags
      }
    ) {
      id
      recordType
      title
      description
      at
      assetChangeIncome {
        asset {
          id
          name
        }
        amount
      }
      assetChangeExpense {
        asset {
          id
          name
        }
        amount
      }
    }
  }
`);

// レコード削除のためのミューテーション
const DeleteRecordDocument = graphql(/* GraphQL */ `
  mutation DeleteRecord($id: ID!) {
    deleteRecord(id: $id) {
      id
    }
  }
`);

export default function RecordDetailPage() {
  const { recordId } = useParams<{ recordId: string }>();
  const navigate = useNavigate();

  // レコードフォームの状態
  const [recordForm, setRecordForm] = useState({
    id: "",
    title: "",
    description: "",
    at: new Date(),
    recordType: "EXPENSE",
    assetID: "",
    fromAssetID: "",
    toAssetID: "",
    amount: 0,
    tags: [] as string[],
    tagInput: "" // カンマ区切りのタグ入力用
  });

  // ダイアログの状態
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  
  // エラー通知の状態
  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");

  // GraphQL クエリの実行
  const [{ data, fetching, error }] = useQuery({
    query: GetRecordDocument,
    variables: { id: recordId || "" }
  });
  
  // ミューテーション関数の取得
  const [, updateIncomeRecord] = useMutation(UpdateIncomeRecordDocument);
  const [, updateExpenseRecord] = useMutation(UpdateExpenseRecordDocument);
  const [, updateTransferRecord] = useMutation(UpdateTransferRecordDocument);
  const [, deleteRecordMutation] = useMutation(DeleteRecordDocument);

  // 編集ダイアログを開く
  const handleOpenEditDialog = () => {
    if (data?.record) {
      // レコードタイプに応じて適切な初期化を行う
      const tags = data.record.tags ? data.record.tags.map((tag: { id: string, name: string }) => tag.name) : [];
      const formData = {
        id: data.record.id,
        title: data.record.title,
        description: data.record.description || "",
        at: new Date(data.record.at),
        recordType: data.record.recordType,
        assetID: "",
        fromAssetID: "",
        toAssetID: "",
        amount: 0,
        tags: tags,
        tagInput: tags.join(', ')
      };

      // レコードタイプに応じて資産とアマウントを設定
      if (data.record.recordType === "INCOME" && data.record.assetChangeIncome) {
        formData.assetID = data.record.assetChangeIncome.asset.id;
        formData.amount = data.record.assetChangeIncome.amount;
      } else if (data.record.recordType === "EXPENSE" && data.record.assetChangeExpense) {
        formData.assetID = data.record.assetChangeExpense.asset.id;
        formData.amount = Math.abs(data.record.assetChangeExpense.amount);
      } else if (data.record.recordType === "TRANSFER") {
        if (data.record.assetChangeExpense) {
          formData.fromAssetID = data.record.assetChangeExpense.asset.id;
        }
        if (data.record.assetChangeIncome) {
          formData.toAssetID = data.record.assetChangeIncome.asset.id;
          formData.amount = data.record.assetChangeIncome.amount;
        }
      }
      
      setRecordForm(formData);
      setEditDialogOpen(true);
    }
  };

  // 削除ダイアログを開く
  const handleOpenDeleteDialog = () => {
    setDeleteDialogOpen(true);
  };

  // フォームの変更を処理
  const handleFormChange = (field: string, value: any) => {
    setRecordForm(prev => ({ ...prev, [field]: value }));
  };

  // レコードを更新
  const handleUpdateRecord = async () => {
    try {
      let result;
      
      // レコードタイプに応じて適切なミューテーションを呼び出す
      if (recordForm.recordType === "INCOME") {
        result = await updateIncomeRecord({
          id: recordForm.id,
          title: recordForm.title,
          description: recordForm.description || "",
          at: recordForm.at.toISOString(),
          assetID: recordForm.assetID,
          amount: recordForm.amount,
          tags: recordForm.tags
        });
      } else if (recordForm.recordType === "EXPENSE") {
        result = await updateExpenseRecord({
          id: recordForm.id,
          title: recordForm.title,
          description: recordForm.description || "",
          at: recordForm.at.toISOString(),
          assetID: recordForm.assetID,
          amount: recordForm.amount,
          tags: recordForm.tags
        });
      } else if (recordForm.recordType === "TRANSFER") {
        result = await updateTransferRecord({
          id: recordForm.id,
          title: recordForm.title,
          description: recordForm.description || "",
          at: recordForm.at.toISOString(),
          fromAssetID: recordForm.fromAssetID,
          toAssetID: recordForm.toAssetID,
          amount: recordForm.amount,
          tags: recordForm.tags
        });
      }
      
      if (result?.error) {
        setErrorMessage(result.error.message);
      } else {
        setSuccessMessage("レコードが正常に更新されました");
        setEditDialogOpen(false);
        // 必要に応じてクエリを再実行
      }
    } catch (err) {
      setErrorMessage("レコードの更新中にエラーが発生しました");
    }
  };

  // レコードを削除
  const handleDeleteRecord = async () => {
    if (!recordId) return;
    
    try {
      const result = await deleteRecordMutation({ id: recordId });
      if (result.error) {
        setErrorMessage(result.error.message);
      } else {
        setSuccessMessage("レコードが正常に削除されました");
        setDeleteDialogOpen(false);
        // 前の画面に戻る
        navigate(-1);
      }
    } catch (err) {
      setErrorMessage("レコードの削除中にエラーが発生しました");
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

  const record = data?.record;
  // タグ関連の機能は削除したため、この変数も不要になりました

  if (!record) {
    return (
      <Paper elevation={3} sx={{ p: 3, maxWidth: 600, mx: "auto", textAlign: "center" }}>
        <Typography variant="h6" color="error">レコードが見つかりません</Typography>
        <Button onClick={() => navigate(-1)} variant="contained" color="primary" sx={{ mt: 2 }}>
          前のページに戻る
        </Button>
      </Paper>
    );
  }

  // レコードタイプに基づいて色を取得する関数
  const getRecordTypeColor = (type: string) => {
    switch (type) {
      case 'EXPENSE': return 'error';
      case 'INCOME': return 'success';
      case 'TRANSFER': return 'info';
      default: return 'default';
    }
  };

  // レコードタイプの日本語表示
  const getRecordTypeLabel = (type: string) => {
    switch (type) {
      case 'EXPENSE': return '支出';
      case 'INCOME': return '収入';
      case 'TRANSFER': return '振替';
      default: return type;
    }
  };

  // 日付をフォーマット
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('ja-JP', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    }).format(date);
  };

  return (
    <Box p={3}>
      {/* ヘッダー */}
      <Stack direction="row" alignItems="center" justifyContent="space-between" mb={3}>
        <Button onClick={() => navigate(-1)} startIcon={<ArrowBackIcon />}>
          前のページに戻る
        </Button>
        <Typography variant="h5">レコード詳細</Typography>
        <Stack direction="row" spacing={1}>
          <Button 
            variant="outlined" 
            color="primary" 
            startIcon={<EditIcon />}
            onClick={handleOpenEditDialog}
          >
            編集
          </Button>
          <Button 
            variant="outlined" 
            color="error" 
            startIcon={<DeleteIcon />}
            onClick={handleOpenDeleteDialog}
          >
            削除
          </Button>
        </Stack>
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

      {/* レコード詳細表示 */}
      <Paper elevation={2} sx={{ p: 3, mb: 3 }}>
        <Box sx={{ mb: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="h4">{record?.title}</Typography>
          <Chip 
            label={getRecordTypeLabel(record?.recordType || "")} 
            color={getRecordTypeColor(record?.recordType || "") as any}
            size="medium"
          />
        </Box>

        <Typography variant="body2" color="text.secondary" gutterBottom>
          {formatDate(record?.at)}
        </Typography>
        
        <Divider sx={{ my: 2 }} />
        
        {record?.recordType === 'INCOME' && record?.assetChangeIncome && (
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, gap: 3 }}>
            <Box>
              <Typography variant="subtitle1">金額</Typography>
              <Typography variant="h5" color="success.main">
                ¥{record.assetChangeIncome.amount.toLocaleString()}
              </Typography>
            </Box>
            
            <Box>
              <Typography variant="subtitle1">入金先</Typography>
              <Typography variant="body1">
                {record.assetChangeIncome.asset.name}
              </Typography>
            </Box>
          </Box>
        )}
        
        {record?.recordType === 'EXPENSE' && record?.assetChangeExpense && (
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, gap: 3 }}>
            <Box>
              <Typography variant="subtitle1">金額</Typography>
              <Typography variant="h5" color="error.main">
                -¥{Math.abs(record.assetChangeExpense.amount).toLocaleString()}
              </Typography>
            </Box>
            
            <Box>
              <Typography variant="subtitle1">支払元</Typography>
              <Typography variant="body1">
                {record.assetChangeExpense.asset.name}
              </Typography>
            </Box>
          </Box>
        )}
        
        {record?.recordType === 'TRANSFER' && record?.assetChangeIncome && record?.assetChangeExpense && (
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, gap: 3 }}>
            <Box>
              <Typography variant="subtitle1">金額</Typography>
              <Typography variant="h5" color="info.main">
                ¥{record.assetChangeIncome.amount.toLocaleString()}
              </Typography>
            </Box>
            
            <Box>
              <Typography variant="subtitle1">振替情報</Typography>
              <Typography variant="body1">
                送金元: {record.assetChangeExpense.asset.name}
              </Typography>
              <Typography variant="body1">
                送金先: {record.assetChangeIncome.asset.name}
              </Typography>
            </Box>
          </Box>
        )}

        {record?.description && (
          <Box sx={{ mt: 3 }}>
            <Typography variant="subtitle1">説明</Typography>
            <Paper variant="outlined" sx={{ p: 2, backgroundColor: 'background.default' }}>
              <Typography variant="body1">{record.description}</Typography>
            </Paper>
          </Box>
        )}

        {/* タグ表示 */}
        {record?.tags && record.tags.length > 0 && (
          <Box sx={{ mt: 3 }}>
            <Typography variant="subtitle1">タグ</Typography>
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mt: 1 }}>
              {record.tags.map((tag: { id: string, name: string }) => (
                <Chip key={tag.id} label={tag.name} size="small" />
              ))}
            </Box>
          </Box>
        )}
      </Paper>

      {/* レコード編集ダイアログ */}
      <LocalizationProvider dateAdapter={AdapterDateFns} adapterLocale={ja}>
        <Dialog open={editDialogOpen} onClose={() => setEditDialogOpen(false)} maxWidth="sm" fullWidth>
          <DialogTitle>レコードを編集</DialogTitle>
          <DialogContent>
            <Box sx={{ mt: 2 }}>
              <Typography variant="subtitle1" gutterBottom>
                レコードタイプ: {record && getRecordTypeLabel(record.recordType)}
              </Typography>

              <TextField
                margin="normal"
                fullWidth
                label="タイトル"
                value={recordForm.title}
                onChange={(e) => handleFormChange('title', e.target.value)}
              />
              
              <TextField
                margin="normal"
                fullWidth
                label="説明"
                multiline
                rows={2}
                value={recordForm.description}
                onChange={(e) => handleFormChange('description', e.target.value)}
              />
              
              <DateTimePicker
                label="日時"
                value={recordForm.at}
                onChange={(date: Date | null) => date && handleFormChange('at', date)}
                sx={{ mt: 2, width: '100%' }}
                format="yyyy年MM月dd日 HH時mm分ss秒"
              />
              
              {record?.recordType === "INCOME" && (
                <Box sx={{ mt: 2 }}>
                  <Typography variant="subtitle2">収入情報</Typography>
                  <TextField
                    margin="normal"
                    fullWidth
                    label="金額"
                    type="number"
                    value={recordForm.amount}
                    onChange={(e) => handleFormChange('amount', parseInt(e.target.value))}
                    InputProps={{
                      startAdornment: <Typography sx={{ mr: 1 }}>¥</Typography>,
                    }}
                  />
                  <Typography variant="body2" sx={{ mt: 1 }}>
                    資産: {record?.assetChangeIncome?.asset.name}
                  </Typography>
                </Box>
              )}
              
              {record?.recordType === "EXPENSE" && (
                <Box sx={{ mt: 2 }}>
                  <Typography variant="subtitle2">支出情報</Typography>
                  <TextField
                    margin="normal"
                    fullWidth
                    label="金額"
                    type="number"
                    value={recordForm.amount}
                    onChange={(e) => handleFormChange('amount', parseInt(e.target.value))}
                    InputProps={{
                      startAdornment: <Typography sx={{ mr: 1 }}>¥</Typography>,
                    }}
                  />
                  <Typography variant="body2" sx={{ mt: 1 }}>
                    資産: {record?.assetChangeExpense?.asset.name}
                  </Typography>
                </Box>
              )}
              
              {record?.recordType === "TRANSFER" && (
                <Box sx={{ mt: 2 }}>
                  <Typography variant="subtitle2">振替情報</Typography>
                  <TextField
                    margin="normal"
                    fullWidth
                    label="金額"
                    type="number"
                    value={recordForm.amount}
                    onChange={(e) => handleFormChange('amount', parseInt(e.target.value))}
                    InputProps={{
                      startAdornment: <Typography sx={{ mr: 1 }}>¥</Typography>,
                    }}
                  />
                  <Typography variant="body2" sx={{ mt: 1 }}>
                    送金元: {record?.assetChangeExpense?.asset.name}
                  </Typography>
                  <Typography variant="body2">
                    送金先: {record?.assetChangeIncome?.asset.name}
                  </Typography>
                </Box>
              )}
              
              {/* タグ編集 */}
              <Box sx={{ mt: 3 }}>
                <Typography variant="subtitle2">タグ</Typography>
                <TextField
                  margin="normal"
                  fullWidth
                  label="タグ (カンマ区切り)"
                  value={recordForm.tagInput}
                  onChange={(e) => {
                    const tagInput = e.target.value;
                    // 入力値をそのまま保存
                    handleFormChange('tagInput', tagInput);
                    
                    // カンマ区切りで分割してタグ配列に変換
                    const tagNames = tagInput
                      .split(',')
                      .map(tag => tag.trim())
                      .filter(tag => tag !== '');
                    
                    handleFormChange('tags', tagNames);
                  }}
                  placeholder="例: 食費, 日用品, 交通費"
                  helperText="カンマ(,)で区切って複数のタグを入力できます"
                />
                {recordForm.tags.length > 0 && (
                  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mt: 1 }}>
                    {recordForm.tags.map((tagName: string) => (
                      <Chip 
                        key={tagName} 
                        label={tagName} 
                        size="small"
                        onDelete={() => {
                          const newTagNames = recordForm.tags.filter(t => t !== tagName);
                          handleFormChange('tags', newTagNames);
                          // tagInputも更新
                          handleFormChange('tagInput', newTagNames.join(', '));
                        }}
                      />
                    ))}
                  </Box>
                )}
              </Box>
            </Box>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setEditDialogOpen(false)}>キャンセル</Button>
            <Button 
              onClick={handleUpdateRecord} 
              variant="contained" 
              color="primary"
              disabled={!recordForm.title.trim() || recordForm.amount <= 0}
            >
              更新
            </Button>
          </DialogActions>
        </Dialog>
      </LocalizationProvider>

      {/* レコード削除確認ダイアログ */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>レコードを削除</DialogTitle>
        <DialogContent>
          <Typography>
            レコード「{record?.title}」を本当に削除しますか？
            この操作は元に戻せません。
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>キャンセル</Button>
          <Button 
            onClick={handleDeleteRecord} 
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
