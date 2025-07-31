import { useState, useMemo } from "react";
import { useParams, Link as RouterLink } from "react-router-dom";
import { useQuery, useMutation } from "urql";
import { graphql } from "../gql";
import type { RecordType } from "../gql/graphql";

// 日時を「yyyy年MM月dd日 HH時mm分ss秒」形式でフォーマットする関数
const formatDateTime = (dateString: string) => {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('ja-JP', {
    year: 'numeric', month: 'long', day: 'numeric',
    hour: '2-digit', minute: '2-digit', second: '2-digit'
  }).format(date);
};
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
  MenuItem,
  FormControl,
  InputLabel,
  Select,
  Chip,
  FormControlLabel,
  Checkbox,
  Card,
  CardContent,
  CardActions
} from "@mui/material";
import { DateTimePicker } from "@mui/x-date-pickers/DateTimePicker";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFnsV3";
import { ja } from "date-fns/locale";
import AddIcon from "@mui/icons-material/Add";
import FilterListIcon from "@mui/icons-material/FilterList";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import ArrowForwardIcon from "@mui/icons-material/ArrowForward";

// GraphQL クエリ定義
const GetMonthlyRecordsDocument = graphql(/* GraphQL */ `
  query GetMonthlyRecords($year: Int!, $month: Int!, $tagNames: [String!], $assetIds: [ID!], $recordTypes: [RecordType!]) {
    # records(year: $year, month: $month, tagNames: $tagNames, assetIds: $assetIds, recordTypes: $recordTypes) {
    recordsPerMonth(first: 100, year: $year, month: $month, tagNames: $tagNames, assetIds: $assetIds, recordTypes: $recordTypes) {
        nodes {
            id
            recordType
            title
            description
            at
            assetChangeIncome { 
                asset { id name }
                amount
            }
            assetChangeExpense {
                asset { id name }
                amount
            }
            tags {
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

        totalAssets
    }

    assets(first: 10000) {
        nodes {
            id
            name
            category {
                id
                name
            }
        }
    }
    
    tags(first: 10000) {
        nodes {
            id
            name
        }
    }
  }
`);

// 新規レコード作成のためのミューテーション
// 収入レコード作成
const CreateIncomeRecordDocument = graphql(/* GraphQL */ `
  mutation CreateIncomeRecord($input: createIncomeRecordInput!) {
    createIncomeRecord(input: $input) {
      id
      recordType
      title
      description
      at
      assetChangeIncome {
        asset { id name }
        amount
      }
    }
  }
`);

// 支出レコード作成
const CreateExpenseRecordDocument = graphql(/* GraphQL */ `
  mutation CreateExpenseRecord($input: createExpenseRecordInput!) {
    createExpenseRecord(input: $input) {
      id
      recordType
      title
      description
      at
      assetChangeExpense {
        asset { id name }
        amount
      }
    }
  }
`);

// 振替レコード作成
const CreateTransferRecordDocument = graphql(/* GraphQL */ `
  mutation CreateTransferRecord($input: createTransferRecordInput!) {
    createTransferRecord(input: $input) {
      id
      recordType
      title
      description
      at
      assetChangeIncome {
        asset { id name }
        amount
      }
      assetChangeExpense {
        asset { id name }
        amount
      }
    }
  }
`)

export default function MonthlyRecordPage() {
  const { year: yearParam, month: monthParam } = useParams<{ year: string, month: string }>();
  const year = parseInt(yearParam || "0", 10);
  const month = parseInt(monthParam || "0", 10);

  // フィルター状態の管理
  const [selectedTags, setSelectedTags] = useState<string[]>([]); // タグ名の配列
  const [selectedAssets, setSelectedAssets] = useState<string[]>([]); // 資産IDの配列
  const [recordTypeFilter, setRecordTypeFilter] = useState<RecordType[]>([]);

  // フィルターダイアログの一時的な状態 (ダイアログ内での選択を保持)
  const [tempSelectedTags, setTempSelectedTags] = useState<string[]>([]);
  const [tempSelectedAssets, setTempSelectedAssets] = useState<string[]>([]);
  const [tempRecordTypeFilter, setTempRecordTypeFilter] = useState<RecordType[]>([]);

  // フィルターダイアログの状態
  const [filterDialogOpen, setFilterDialogOpen] = useState(false);

  // 新規作成ダイアログの状態
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [newRecord, setNewRecord] = useState({
    title: "",
    description: "",
    recordType: "EXPENSE" as RecordType,
    at: new Date(),
    assetId: "",
    fromAssetId: "",
    toAssetId: "",
    amount: 0,
    tagIds: [] as string[],
    tagNames: [] as string[], // タグ名の配列を追加
    tagInput: "" // カンマを含むタグ入力用の文字列
  });
  
  // ミューテーション関数（codegenの実行後にエラーは解消されるため、一時的にignoreします）
  // @ts-ignore
  const [, createIncomeRecord] = useMutation(CreateIncomeRecordDocument);
  // @ts-ignore
  const [, createExpenseRecord] = useMutation(CreateExpenseRecordDocument);
  // @ts-ignore
  const [, createTransferRecord] = useMutation(CreateTransferRecordDocument);

  // GraphQLクエリの実行
  // codegenの実行後にエラーは解消されるため、一時的にignoreします
  // @ts-ignore
  const [{ data, fetching, error }] = useQuery({
    query: GetMonthlyRecordsDocument,
    variables: {
      year,
      month,
      tagNames: selectedTags.length > 0 ? selectedTags : undefined, // タグ名をそのまま送信
      assetIds: selectedAssets.length > 0 ? selectedAssets : undefined,
      recordTypes: recordTypeFilter.length > 0 ? recordTypeFilter : undefined,
    }
  });

  // 月のナビゲーション用の計算
  const prevMonth = useMemo(() => {
    const date = new Date(year, month - 1);
    date.setMonth(date.getMonth() - 1);
    return {
      year: date.getFullYear(),
      month: date.getMonth() + 1
    };
  }, [year, month]);

  const nextMonth = useMemo(() => {
    const date = new Date(year, month - 1);
    date.setMonth(date.getMonth() + 1);
    return {
      year: date.getFullYear(),
      month: date.getMonth() + 1
    };
  }, [year, month]);

  // フィルターダイアログを開く
  const handleOpenFilterDialog = () => {
    // 現在のフィルター状態を一時的な状態にコピー
    setTempSelectedTags([...selectedTags]);
    setTempSelectedAssets([...selectedAssets]);
    setTempRecordTypeFilter([...recordTypeFilter]);
    setFilterDialogOpen(true);
  };

  // フィルターダイアログを閉じる (キャンセル)
  const handleCloseFilterDialog = () => {
    setFilterDialogOpen(false);
  };

  // フィルターを適用する
  const handleApplyFilter = () => {
    // 一時的な状態を実際のフィルター状態に反映
    setSelectedTags(tempSelectedTags);
    setSelectedAssets(tempSelectedAssets);
    setRecordTypeFilter(tempRecordTypeFilter);
    setFilterDialogOpen(false);
  };

  // フィルターをリセットする
  const handleResetFilter = () => {
    // 一時的な状態をクリア
    setTempSelectedTags([]);
    setTempSelectedAssets([]);
    setTempRecordTypeFilter([]);
    
    // 実際のフィルター状態もクリア
    setSelectedTags([]);
    setSelectedAssets([]);
    setRecordTypeFilter([]);
    setFilterDialogOpen(false);
  };

  // 新規レコード作成ダイアログを開く
  const handleOpenCreateDialog = () => {
    setCreateDialogOpen(true);
  };

  // 新規レコード作成ダイアログを閉じる
  const handleCloseCreateDialog = () => {
    setCreateDialogOpen(false);
    // フォームをリセット
    setNewRecord({
      title: "",
      description: "",
      recordType: "EXPENSE" as RecordType,
      at: new Date(),
      assetId: "",
      fromAssetId: "",
      toAssetId: "",
      amount: 0,
      tagIds: [],
      tagNames: [],
      tagInput: ""
    });
  };

  // 新規レコード作成フォームの変更を処理する
  const handleNewRecordChange = (field: string, value: any) => {
    setNewRecord(prev => ({ ...prev, [field]: value }));
  };

  // 新規レコードを保存する
  const handleSaveRecord = async () => {
    try {
      let result;

      // レコードタイプに応じて適切なミューテーションを実行
      switch (newRecord.recordType) {
        case "INCOME":
          result = await createIncomeRecord({
            input: {
              title: newRecord.title,
              description: newRecord.description,
              at: newRecord.at.toISOString(),
              assetID: newRecord.assetId,
              amount: newRecord.amount,
              tags: newRecord.tagNames // タグはname配列で送信
            }
          });
          break;
        
        case "EXPENSE":
          result = await createExpenseRecord({
            input: {
              title: newRecord.title,
              description: newRecord.description,
              at: newRecord.at.toISOString(),
              assetID: newRecord.assetId,
              amount: newRecord.amount,
              tags: newRecord.tagNames // タグはname配列で送信
            }
          });
          break;
        
        case "TRANSFER":
          result = await createTransferRecord({
            input: {
              title: newRecord.title,
              description: newRecord.description,
              at: newRecord.at.toISOString(),
              fromAssetID: newRecord.fromAssetId,
              toAssetID: newRecord.toAssetId,
              amount: newRecord.amount,
              tags: newRecord.tagNames // タグはname配列で送信
            }
          });
          break;
      }
      
      if (result?.error) {
        // エラー処理
        console.error('レコード作成エラー:', result.error);
      } else {
        // 成功した場合、クエリを再実行して最新のレコードを取得
        console.log('レコードが作成されました');
      }
      
      handleCloseCreateDialog();
    } catch (err) {
      console.error('レコード作成中にエラーが発生しました:', err);
    }
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

  const records = data?.recordsPerMonth?.nodes || [];
  const initialTotalAssets = data?.recordsPerMonth?.totalAssets || 0;
  
  // 日付でレコードをソート
  const sortedRecords = [...records].sort((a, b) => new Date(a.at).getTime() - new Date(b.at).getTime());
  
  const assets = data?.assets?.nodes || [];
  const tags = data?.tags?.nodes || [];

  // レコードタイプに基づいて色を取得する関数
  const getRecordTypeColor = (type: RecordType) => {
    switch (type) {
      case "EXPENSE": return 'error';
      case "INCOME": return 'success';
      case "TRANSFER": return 'info';
      default: return 'default';
    }
  };

  // レコードタイプの日本語表示
  const getRecordTypeLabel = (type: RecordType) => {
    switch (type) {
      case "EXPENSE": return '支出';
      case "INCOME": return '収入';
      case "TRANSFER": return '振替';
      default: return type;
    }
  };

  // 月の合計を計算
  const totalExpense = records
    .filter((r) => r.recordType === "EXPENSE")
    .reduce((sum: number, record) => {
      const amount = record.assetChangeExpense?.amount || 0;
      return sum + Math.abs(amount);
    }, 0);

  const totalIncome = records
    .filter((r) => r.recordType === "INCOME")
    .reduce((sum: number, record) => {
      const amount = record.assetChangeIncome?.amount || 0;
      return sum + amount;
    }, 0);
    
  // 各レコードでの資産変化と累積総資産を計算する関数
  const calculateCumulativeAssets = () => {
    let cumulativeAssets = initialTotalAssets;
    const recordsWithCumulativeAssets = sortedRecords.map((record) => {
      let assetChange = 0;
      
      if (record.recordType === "INCOME" && record.assetChangeIncome) {
        assetChange = record.assetChangeIncome.amount;
      } else if (record.recordType === "EXPENSE" && record.assetChangeExpense) {
        assetChange = -Math.abs(record.assetChangeExpense.amount);
      }
      // 振替の場合、総資産には影響なし
      
      cumulativeAssets += assetChange;
      return {
        ...record,
        assetChange,
        cumulativeAssets
      };
    });
    
    return recordsWithCumulativeAssets;
  };
  
  // 累積総資産付きのレコード配列
  const recordsWithAssets = calculateCumulativeAssets();

  return (
    <Box p={3}>
      {/* ヘッダー */}
      <Stack direction="row" alignItems="center" justifyContent="space-between" mb={3}>
        <Button component={RouterLink} to="/" startIcon={<ArrowBackIcon />}>
          ホームに戻る
        </Button>
        <Typography variant="h5">{year}年{month}月のレコード</Typography>
        <Box>
          <Button 
            variant="outlined" 
            startIcon={<FilterListIcon />} 
            onClick={handleOpenFilterDialog}
            sx={{ mr: 1 }}
          >
            フィルター
          </Button>
          <Button 
            variant="contained" 
            color="primary" 
            startIcon={<AddIcon />}
            onClick={handleOpenCreateDialog}
          >
            新規作成
          </Button>
        </Box>
      </Stack>

      {/* 月ナビゲーション */}
      <Stack 
        direction="row" 
        spacing={2} 
        justifyContent="center" 
        alignItems="center"
        mb={3}
      >
        <Button
          component={RouterLink}
          to={`/record/monthly/${prevMonth.year}/${prevMonth.month}`}
          startIcon={<ArrowBackIcon />}
        >
          前月
        </Button>
        <Typography variant="h6">
          {year}年{month}月
        </Typography>
        <Button
          component={RouterLink}
          to={`/record/monthly/${nextMonth.year}/${nextMonth.month}`}
          endIcon={<ArrowForwardIcon />}
        >
          次月
        </Button>
      </Stack>

      {/* 月次サマリー */}
      <Paper elevation={2} sx={{ p: 2, mb: 3 }}>
        <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: 'repeat(4, 1fr)' }, gap: 3 }}>
          <Box>
            <Typography variant="subtitle1" gutterBottom>収入合計</Typography>
            <Typography variant="h5" color="success.main">¥{totalIncome.toLocaleString()}</Typography>
          </Box>
          <Box>
            <Typography variant="subtitle1" gutterBottom>支出合計</Typography>
            <Typography variant="h5" color="error.main">¥{totalExpense.toLocaleString()}</Typography>
          </Box>
          <Box>
            <Typography variant="subtitle1" gutterBottom>収支</Typography>
            <Typography variant="h5" color={(totalIncome - totalExpense) >= 0 ? "success.main" : "error.main"}>
              ¥{(totalIncome - totalExpense).toLocaleString()}
            </Typography>
          </Box>
          <Box>
            <Typography variant="subtitle1" gutterBottom>現在の総資産</Typography>
            <Typography variant="h5" color="info.main">
              ¥{recordsWithAssets.length > 0 
                ? recordsWithAssets[recordsWithAssets.length - 1].cumulativeAssets.toLocaleString()
                : initialTotalAssets.toLocaleString()}
            </Typography>
          </Box>
        </Box>
      </Paper>

      {/* フィルター表示 */}
      {(selectedTags.length > 0 || selectedAssets.length > 0 || recordTypeFilter.length > 0) && (
        <Box sx={{ mb: 2 }}>
          <Stack direction="row" spacing={1} flexWrap="wrap">
            {selectedTags.map(tagName => (
              <Chip 
                key={tagName} 
                label={tagName} 
                color="primary" 
                variant="outlined" 
                onDelete={() => setSelectedTags(prev => prev.filter(name => name !== tagName))} 
              />
            ))}
            {selectedAssets.map(assetId => {
              const asset = assets.find(a => a.id === assetId);
              return asset ? (
                <Chip 
                  key={asset.id} 
                  label={asset.name} 
                  color="success" 
                  variant="outlined" 
                  onDelete={() => setSelectedAssets(prev => prev.filter(id => id !== assetId))} 
                />
              ) : null;
            })}
            {recordTypeFilter.map(type => (
              <Chip 
                key={type} 
                label={getRecordTypeLabel(type)} 
                color={getRecordTypeColor(type) as any} 
                variant="outlined" 
                onDelete={() => setRecordTypeFilter(prev => prev.filter(t => t !== type))} 
              />
            ))}
            <Chip 
              label="フィルターをクリア" 
              variant="outlined" 
              onClick={handleResetFilter} 
            />
          </Stack>
        </Box>
      )}

      {/* レコードリスト */}
      {recordsWithAssets.length === 0 ? (
        <Paper elevation={1} sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="h6" color="text.secondary">
            レコードがありません
          </Typography>
          <Button 
            variant="contained" 
            color="primary"
            startIcon={<AddIcon />} 
            onClick={handleOpenCreateDialog}
            sx={{ mt: 2 }}
          >
            新規レコードを作成
          </Button>
        </Paper>
      ) : (
        <Stack spacing={2}>
          {recordsWithAssets.map((record) => {
            let amount = 0;
            let assetName = '';
            
            // レコードタイプに応じた金額と資産名の取得
            if (record.recordType === "INCOME" && record.assetChangeIncome) {
              amount = record.assetChangeIncome.amount;
              assetName = record.assetChangeIncome.asset.name;
            } else if (record.recordType === "EXPENSE" && record.assetChangeExpense) {
              amount = -Math.abs(record.assetChangeExpense.amount);
              assetName = record.assetChangeExpense.asset.name;
            } else if (record.recordType === "TRANSFER") {
              // 振替の場合は入金と出金の両方の情報を表示するための準備
              const fromAsset = record.assetChangeExpense?.asset.name || '';
              const toAsset = record.assetChangeIncome?.asset.name || '';
              assetName = `${fromAsset} → ${toAsset}`;
              amount = record.assetChangeExpense?.amount || 0;
            }

            return (
              <Card key={record.id} sx={{ mb: 1 }}>
                <CardContent>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                    <Typography variant="h6">{record.title}</Typography>
                    <Chip 
                      label={getRecordTypeLabel(record.recordType)} 
                      color={getRecordTypeColor(record.recordType) as any}
                      size="small" 
                    />
                  </Box>

                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    {formatDateTime(record.at)}
                  </Typography>
                  
                  <Typography variant="h6" color={amount < 0 ? "error.main" : "success.main"}>
                    {amount < 0 ? '-' : ''}¥{Math.abs(amount).toLocaleString()}
                  </Typography>
                  
                  <Typography variant="body2">
                    {assetName}
                  </Typography>
                  
                  {/* 総資産表示 */}
                  <Box sx={{ mt: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Typography variant="body2" color="text.secondary">
                      この取引後の総資産:
                    </Typography>
                    <Typography variant="body1" fontWeight="bold">
                      ¥{record.cumulativeAssets.toLocaleString()}
                    </Typography>
                  </Box>

                  {record.description && (
                    <Typography variant="body2" sx={{ mt: 1 }}>
                      {record.description}
                    </Typography>
                  )}

                  {record.tags && record.tags.length > 0 && (
                    <Box sx={{ mt: 1, display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                      {record.tags.map((tag: { id: string, name: string }) => (
                        <Chip key={tag.id} label={tag.name} size="small" />
                      ))}
                    </Box>
                  )}
                </CardContent>
                <CardActions>
                  <Button 
                    size="small" 
                    component={RouterLink} 
                    to={`/record/${record.id}`}
                  >
                    詳細
                  </Button>
                </CardActions>
              </Card>
            );
          })}
        </Stack>
      )}

      {/* フィルターダイアログ */}
      <Dialog open={filterDialogOpen} onClose={handleCloseFilterDialog} maxWidth="sm" fullWidth>
        <DialogTitle>フィルター設定</DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2 }}>
            <Typography variant="subtitle1" gutterBottom>レコードタイプ</Typography>
            <FormControlLabel
              control={
                <Checkbox 
                  checked={tempRecordTypeFilter.includes("EXPENSE")}
                  onChange={(e) => {
                    if (e.target.checked) {
                      setTempRecordTypeFilter(prev => [...prev, "EXPENSE"]);
                    } else {
                      setTempRecordTypeFilter(prev => prev.filter(t => t !== "EXPENSE"));
                    }
                  }}
                />
              }
              label="支出"
            />
            <FormControlLabel
              control={
                <Checkbox 
                  checked={tempRecordTypeFilter.includes("INCOME")}
                  onChange={(e) => {
                    if (e.target.checked) {
                      setTempRecordTypeFilter(prev => [...prev, "INCOME"]);
                    } else {
                      setTempRecordTypeFilter(prev => prev.filter(t => t !== "INCOME"));
                    }
                  }}
                />
              }
              label="収入"
            />
            <FormControlLabel
              control={
                <Checkbox 
                  checked={tempRecordTypeFilter.includes("TRANSFER")}
                  onChange={(e) => {
                    if (e.target.checked) {
                      setTempRecordTypeFilter(prev => [...prev, "TRANSFER"]);
                    } else {
                      setTempRecordTypeFilter(prev => prev.filter(t => t !== "TRANSFER"));
                    }
                  }}
                />
              }
              label="振替"
            />
          </Box>
          
          <Box sx={{ mt: 3 }}>
            <Typography variant="subtitle1" gutterBottom>資産</Typography>
            <FormControl fullWidth sx={{ mt: 1 }}>
              <InputLabel>資産で絞り込み</InputLabel>
              <Select
                multiple
                value={tempSelectedAssets}
                onChange={(e) => setTempSelectedAssets(e.target.value as string[])}
                renderValue={(selected) => (
                  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                    {selected.map((assetId) => {
                      const asset = assets.find(a => a.id === assetId);
                      return asset ? (
                        <Chip key={assetId} label={asset.name} size="small" />
                      ) : null;
                    })}
                  </Box>
                )}
              >
                {assets.map(asset => (
                  <MenuItem key={asset.id} value={asset.id}>
                    {asset.name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Box>
          
          <Box sx={{ mt: 3 }}>
            <Typography variant="subtitle1" gutterBottom>タグ</Typography>
            <FormControl fullWidth sx={{ mt: 1 }}>
              <InputLabel>タグで絞り込み</InputLabel>
              <Select
                multiple
                value={tempSelectedTags}
                onChange={(e) => {
                  const tagNames = e.target.value as string[];
                  setTempSelectedTags(tagNames);
                }}
                renderValue={(selected) => (
                  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                    {selected.map((tagName) => (
                      <Chip key={tagName} label={tagName} size="small" />
                    ))}
                  </Box>
                )}
              >
                {tags.map(tag => (
                  <MenuItem key={tag.id} value={tag.name}>
                    {tag.name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleResetFilter}>リセット</Button>
          <Button onClick={handleApplyFilter} variant="contained" color="primary">
            適用
          </Button>
        </DialogActions>
      </Dialog>

      {/* 新規レコード作成ダイアログ */}
      <LocalizationProvider dateAdapter={AdapterDateFns} adapterLocale={ja}>
        <Dialog open={createDialogOpen} onClose={handleCloseCreateDialog} maxWidth="sm" fullWidth>
          <DialogTitle>新規レコード作成</DialogTitle>
          <DialogContent>
            <Box sx={{ mt: 2 }}>
              <FormControl fullWidth margin="normal">
                <InputLabel id="record-type-label">レコードタイプ</InputLabel>
                <Select
                  labelId="record-type-label"
                  value={newRecord.recordType}
                  onChange={(e) => handleNewRecordChange('recordType', e.target.value)}
                  label="レコードタイプ"
                >
                  <MenuItem value="EXPENSE">支出</MenuItem>
                  <MenuItem value="INCOME">収入</MenuItem>
                  <MenuItem value="TRANSFER">振替</MenuItem>
                </Select>
              </FormControl>
              
              <TextField
                margin="normal"
                fullWidth
                label="タイトル"
                value={newRecord.title}
                onChange={(e) => handleNewRecordChange('title', e.target.value)}
              />
              
              <TextField
                margin="normal"
                fullWidth
                label="説明"
                multiline
                rows={2}
                value={newRecord.description}
                onChange={(e) => handleNewRecordChange('description', e.target.value)}
              />
              
              <DateTimePicker
                label="日時"
                value={newRecord.at}
                onChange={(date: Date | null) => date && handleNewRecordChange('at', date)}
                sx={{ mt: 2, width: '100%' }}
                format="yyyy年MM月dd日 HH時mm分ss秒"
              />
              
              <TextField
                margin="normal"
                fullWidth
                label="金額"
                type="number"
                value={newRecord.amount}
                onChange={(e) => handleNewRecordChange('amount', parseInt(e.target.value, 10))}
                InputProps={{
                  startAdornment: <Typography sx={{ mr: 1 }}>¥</Typography>,
                }}
              />
              
              {/* 収入と支出の場合の資産選択 */}
              {(newRecord.recordType === "INCOME" || newRecord.recordType === "EXPENSE") && (
                <FormControl fullWidth margin="normal">
                  <InputLabel id="asset-label">
                    {newRecord.recordType === "INCOME" ? "入金先" : "支払元"}
                  </InputLabel>
                  <Select
                    labelId="asset-label"
                    value={newRecord.assetId}
                    onChange={(e) => handleNewRecordChange('assetId', e.target.value)}
                    label={newRecord.recordType === "INCOME" ? "入金先" : "支払元"}
                  >
                    {assets.map(asset => (
                      <MenuItem key={asset.id} value={asset.id}>
                        {asset.name}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              )}
              
              {/* 振替の場合は送金元と送金先の両方を選択 */}
              {newRecord.recordType === "TRANSFER" && (
                <>
                  <FormControl fullWidth margin="normal">
                    <InputLabel id="from-asset-label">送金元</InputLabel>
                    <Select
                      labelId="from-asset-label"
                      value={newRecord.fromAssetId}
                      onChange={(e) => handleNewRecordChange('fromAssetId', e.target.value)}
                      label="送金元"
                    >
                      {assets.map(asset => (
                        <MenuItem key={asset.id} value={asset.id}>
                          {asset.name}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                  
                  <FormControl fullWidth margin="normal">
                    <InputLabel id="to-asset-label">送金先</InputLabel>
                    <Select
                      labelId="to-asset-label"
                      value={newRecord.toAssetId}
                      onChange={(e) => handleNewRecordChange('toAssetId', e.target.value)}
                      label="送金先"
                    >
                      {assets.map(asset => (
                        <MenuItem key={asset.id} value={asset.id}>
                          {asset.name}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </>
              )}
              
              <TextField
                margin="normal"
                fullWidth
                label="タグ (カンマ区切り)"
                value={newRecord.tagInput || ''}
                onChange={(e) => {
                  const tagInput = e.target.value;
                  // 入力値をそのまま保存
                  handleNewRecordChange('tagInput', tagInput);
                  
                  // カンマ区切りで分割してタグ配列に変換
                  const tagNames = tagInput
                    .split(',')
                    .map(tag => tag.trim())
                    .filter(tag => tag !== '');
                  
                  handleNewRecordChange('tagNames', tagNames);
                  
                  // IDリストも更新（既存コードとの互換性のため）
                  const selectedIds = tagNames.map((name: string) => {
                    const tag = tags.find(t => t.name === name);
                    return tag ? tag.id : '';
                  }).filter(Boolean);
                  handleNewRecordChange('tagIds', selectedIds);
                }}
                placeholder="例: 食費, 日用品, 交通費"
                helperText="カンマ(,)で区切って複数のタグを入力できます"
              />
              
              {newRecord.tagNames.length > 0 && (
                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mt: 1 }}>
                  {newRecord.tagNames.map((tagName: string) => (
                    <Chip 
                      key={tagName} 
                      label={tagName} 
                      size="small"
                      onDelete={() => {
                        const newTagNames = newRecord.tagNames.filter(t => t !== tagName);
                        handleNewRecordChange('tagNames', newTagNames);
                        // tagInputも更新
                        handleNewRecordChange('tagInput', newTagNames.join(', '));
                      }}
                    />
                  ))}
                </Box>
              )}
            </Box>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCloseCreateDialog}>キャンセル</Button>
            <Button 
              onClick={handleSaveRecord} 
              variant="contained" 
              color="primary"
              disabled={
                !newRecord.title.trim() || 
                newRecord.amount <= 0 || 
                (newRecord.recordType === "INCOME" && !newRecord.assetId) || 
                (newRecord.recordType === "EXPENSE" && !newRecord.assetId) || 
                (newRecord.recordType === "TRANSFER" && (!newRecord.fromAssetId || !newRecord.toAssetId))
              }
            >
              保存
            </Button>
          </DialogActions>
        </Dialog>
      </LocalizationProvider>
    </Box>
  );
}
