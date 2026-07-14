'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  App,
  Button,
  Col,
  Descriptions,
  Divider,
  Empty,
  Form,
  Input,
  Pagination,
  Row,
  Select,
  Space,
  Switch,
  Splitter,
  Table,
  Tag,
  Timeline,
  Typography,
} from 'antd';
import type { TableColumnsType } from 'antd';
import { useEffect, useState } from 'react';

import {
  createStockAction,
  listInventoryStocks,
  listDealerships,
  listStockHistory,
} from '@/features/inventory/inventory-service';
import type {
  CreateActionInput,
  InventoryFilters,
  InventoryActionType,
  InventoryStock,
  StockHistoryEvent,
} from '@/features/inventory/types';

const currencyFormatter = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
  maximumFractionDigits: 0,
});

type FilterFormValues = {
  search?: string;
  ageBand?: 'all' | '0-30' | '31-90' | '91+';
  agingOnly?: boolean;
};

const actionOptions: Array<{ label: string; value: InventoryActionType }> = [
  { label: 'Price reduction planned', value: 'PRICE_REDUCTION_PLANNED' },
  { label: 'Transfer proposed', value: 'TRANSFER_PROPOSED' },
  { label: 'Marketing campaign', value: 'MARKETING_CAMPAIGN' },
  { label: 'Awaiting review', value: 'AWAITING_REVIEW' },
  { label: 'Other', value: 'OTHER' },
];

const actionLabels: Record<InventoryActionType, string> = {
  PRICE_REDUCTION_PLANNED: 'Price reduction planned',
  TRANSFER_PROPOSED: 'Transfer proposed',
  MARKETING_CAMPAIGN: 'Marketing campaign',
  AWAITING_REVIEW: 'Awaiting review',
  OTHER: 'Other',
};

function toInventoryFilters(values: FilterFormValues): InventoryFilters {
  const baseFilters: InventoryFilters = {
    search: values.search,
    agingOnly: values.agingOnly,
  };
  switch (values.ageBand) {
    case '0-30':
      return { ...baseFilters, minAgeDays: 0, maxAgeDays: 30 };
    case '31-90':
      return { ...baseFilters, minAgeDays: 31, maxAgeDays: 90 };
    case '91+':
      return { ...baseFilters, minAgeDays: 91 };
    default:
      return baseFilters;
  }
}

export function InventoryDashboard() {
  const { message } = App.useApp();
  const queryClient = useQueryClient();
  const [filterForm] = Form.useForm<FilterFormValues>();
  const [dealershipId, setDealershipId] = useState('');
  const [filters, setFilters] = useState<InventoryFilters>({});
  const [selectedStockId, setSelectedStockId] = useState<string>();
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 12;

  const dealershipsQuery = useQuery({
    queryKey: ['dealerships'],
    queryFn: listDealerships,
  });

  const dealerships = dealershipsQuery.data;
  const activeDealershipId = dealershipId || dealerships?.[0]?.id || '';

  const inventoryQuery = useQuery({
    queryKey: ['inventory-stocks', activeDealershipId, filters],
    queryFn: () => listInventoryStocks(activeDealershipId, { pageSize: 100, ...filters }),
    enabled: Boolean(activeDealershipId),
  });

  const stocks = inventoryQuery.data?.items ?? [];
  const totalStocks = inventoryQuery.data?.total ?? 0;
  const paginatedStocks = stocks.slice((currentPage - 1) * pageSize, currentPage * pageSize);
  const selectedStock = paginatedStocks.find((stock) => stock.id === selectedStockId) ?? paginatedStocks[0];

  const historyQuery = useQuery({
    queryKey: ['stock-history', selectedStock?.dealershipId, selectedStock?.id],
    queryFn: () => listStockHistory(selectedStock!.dealershipId, selectedStock!.id),
    enabled: Boolean(selectedStock),
  });
  const history = historyQuery.data ?? [];

  const createActionMutation = useMutation({
    mutationFn: (values: CreateActionInput) => {
      if (!selectedStock) throw new Error('No stock selected');
      return createStockAction(activeDealershipId, selectedStock.id, values);
    },
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['inventory-stocks', activeDealershipId] }),
        selectedStock
          ? queryClient.invalidateQueries({
              queryKey: ['stock-history', selectedStock.dealershipId, selectedStock.id],
            })
          : Promise.resolve(),
      ]);
      message.success('Aging-stock action saved.');
    },
    onError: () => {
      message.error('Unable to save the action.');
    },
  });

  async function handleCreateAction(values: CreateActionInput) {
    await createActionMutation.mutateAsync(values);
  }

  const currentDealership = dealerships?.find((item) => item.id === activeDealershipId);
  const loadError = dealershipsQuery.isError || inventoryQuery.isError
    ? 'Unable to load inventory. Check the selected data source.'
    : '';
  const loading = dealershipsQuery.isLoading || inventoryQuery.isLoading || inventoryQuery.isFetching;

  const columns: TableColumnsType<InventoryStock> = [
    {
      title: 'Vehicle',
      key: 'vehicle',
      fixed: 'left',
      width: 210,
      render: (_, stock) => (
        <Space orientation="vertical" size={0}>
          <Typography.Text strong>{stock.make} {stock.model}</Typography.Text>
          <Typography.Text type="secondary">{stock.modelYear} · {stock.vin}</Typography.Text>
        </Space>
      ),
    },
    {
      title: 'Stock age',
      dataIndex: 'inventoryAgeDays',
      key: 'age',
      sorter: (left, right) => left.inventoryAgeDays - right.inventoryAgeDays,
      render: (days: number, stock) => (
        <Space>
          <Typography.Text>{days} days</Typography.Text>
          {stock.isAging && <Tag color="red">Aging</Tag>}
        </Space>
      ),
    },
    {
      title: 'Stocked in',
      dataIndex: 'stockedInAt',
      key: 'stockedInAt',
      render: (value: string) => new Date(value).toLocaleDateString(),
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (value: InventoryStock['status']) => (
        <Tag color={value === 'IN_STOCK' ? 'green' : 'default'}>
          {value === 'IN_STOCK' ? 'In stock' : 'Out of stock'}
        </Tag>
      ),
    },
    {
      title: 'Price',
      dataIndex: 'price',
      key: 'price',
      align: 'right',
      render: (value: number) => currencyFormatter.format(value),
    },
    {
      title: 'Current action',
      key: 'action',
      width: 220,
      render: (_, stock) => stock.latestAction
        ? <Tag color="blue">{actionLabels[stock.latestAction.actionType]}</Tag>
        : <Typography.Text type="secondary">No action logged</Typography.Text>,
    },
    {
      title: '',
      key: 'manage',
      fixed: 'right',
      width: 130,
      render: (_, stock) => (
        <Button
          type={selectedStock?.id === stock.id ? 'primary' : 'default'}
          onClick={() => setSelectedStockId(stock.id)}
        >
          Details
        </Button>
      ),
    },
  ];

  return (
    <main className="mx-auto flex w-full max-w-7xl flex-col gap-3 px-4 py-4 sm:px-6">
      <section className="flex shrink-0 flex-col justify-between gap-3 md:flex-row md:items-end">
        <div>
          <Typography.Title level={2} className="mb-1!">Inventory dashboard</Typography.Title>
          <Typography.Text type="secondary">
            Monitor current stock and act on vehicles held longer than 90 days.
          </Typography.Text>
        </div>
        <div className="w-full md:w-72">
          <Typography.Text className="mb-1 block text-xs" type="secondary">Dealership</Typography.Text>
          <Select
            className="w-full"
            value={activeDealershipId || undefined}
            loading={!(dealerships?.length) && loading}
            options={(dealerships ?? []).map((dealership) => ({
              label: `${dealership.name} · ${dealership.location}`,
              value: dealership.id,
            }))}
            onChange={(value) => {
              setDealershipId(value);
              setSelectedStockId(undefined);
              setCurrentPage(1);
            }}
          />
        </div>
      </section>

      <section className="shrink-0 rounded-lg border border-neutral-200 bg-white px-4 pt-4">
        <Form
          form={filterForm}
          layout="vertical"
          initialValues={{ ageBand: 'all', agingOnly: false }}
          onFinish={(values) => {
            setFilters(toInventoryFilters(values));
            setCurrentPage(1);
          }}
        >
          <Row gutter={16} align="bottom">
            <Col xs={24} sm={12} lg={8}>
              <Form.Item name="search" label="Search">
                <Input allowClear placeholder="Search make or model" />
              </Form.Item>
            </Col>
            <Col xs={24} sm={12} lg={5}>
              <Form.Item name="ageBand" label="Inventory age">
                <Select options={[
                  { label: 'All ages', value: 'all' },
                  { label: '0–30 days', value: '0-30' },
                  { label: '31–90 days', value: '31-90' },
                  { label: '91+ days', value: '91+' },
                ]} />
              </Form.Item>
            </Col>
            <Col xs={24} sm={12} lg={4}>
              <Form.Item name="agingOnly" label="Aging stock only" valuePropName="checked">
                <Switch />
              </Form.Item>
            </Col>
            <Col xs={24} lg={5}>
              <Form.Item>
                <Space>
                  <Button type="primary" htmlType="submit">Apply filters</Button>
                  <Button onClick={() => {
                    filterForm.resetFields();
                    setFilters({});
                    setCurrentPage(1);
                  }}>Reset</Button>
                </Space>
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </section>

      <Splitter className="rounded-lg border border-neutral-200 bg-white">
        <Splitter.Panel defaultSize="68%" min="48%">
          <div className="p-4">
            <div className="mb-4 flex flex-col justify-between gap-3 xl:flex-row xl:items-center">
              <Typography.Title level={4} className="m-0!">
                {currentDealership ? `${currentDealership.name} inventory` : 'Inventory'}
              </Typography.Title>
              <Space className="justify-between xl:justify-end" wrap>
                <Typography.Text type="secondary">{totalStocks} results</Typography.Text>
                <TableTopPagination
                  current={currentPage}
                  pageSize={pageSize}
                  total={totalStocks}
                  onChange={(page) => {
                    setCurrentPage(page);
                    setSelectedStockId(undefined);
                  }}
                />
              </Space>
            </div>
            {loadError ? (
              <Empty description={loadError}>
                <Button onClick={() => void inventoryQuery.refetch()}>Try again</Button>
              </Empty>
            ) : (
              <Table
                rowKey="id"
                size='small'
                loading={loading}
                columns={columns}
                dataSource={paginatedStocks}
                pagination={false}
                rowClassName={(stock) => selectedStock?.id === stock.id ? 'bg-blue-50' : 'cursor-pointer'}
                scroll={{ x: 1160 }}
                onRow={(stock) => ({
                  onClick: () => setSelectedStockId(stock.id),
                })}
              />
            )}
          </div>
        </Splitter.Panel>
        <Splitter.Panel defaultSize="32%" min="28%" collapsible>
          <StockDetailsPanel
            actionLabels={actionLabels}
            actionOptions={actionOptions}
            history={history}
            historyLoading={historyQuery.isLoading || historyQuery.isFetching}
            saving={createActionMutation.isPending}
            selectedStock={selectedStock}
            onClose={() => setSelectedStockId(undefined)}
            onSubmit={handleCreateAction}
          />
        </Splitter.Panel>
      </Splitter>
    </main>
  );
}

type StockDetailsPanelProps = {
  actionLabels: Record<InventoryActionType, string>;
  actionOptions: Array<{ label: string; value: InventoryActionType }>;
  history: StockHistoryEvent[];
  historyLoading: boolean;
  saving: boolean;
  selectedStock?: InventoryStock;
  onClose: () => void;
  onSubmit: (values: CreateActionInput) => Promise<void>;
};

type TableTopPaginationProps = {
  current: number;
  pageSize: number;
  total: number;
  onChange: (page: number) => void;
};

function TableTopPagination({ current, pageSize, total, onChange }: TableTopPaginationProps) {
  if (total <= pageSize) return null;

  return (
    <Pagination
      current={current}
      pageSize={pageSize}
      total={total}
      showSizeChanger={false}
      size="small"
      onChange={onChange}
    />
  );
}

function StockDetailsPanel({
  actionLabels,
  actionOptions,
  history,
  historyLoading,
  saving,
  selectedStock,
  onClose,
  onSubmit,
}: StockDetailsPanelProps) {
  if (!selectedStock) {
    return (
      <div className="flex h-full items-center justify-center p-6">
        <Empty description="Select a stock item to view details" />
      </div>
    );
  }

  return (
    <StockDetailsContent
      actionLabels={actionLabels}
      actionOptions={actionOptions}
      history={history}
      historyLoading={historyLoading}
      saving={saving}
      selectedStock={selectedStock}
      onClose={onClose}
      onSubmit={onSubmit}
    />
  );
}

function StockDetailsContent({
  actionLabels,
  actionOptions,
  history,
  historyLoading,
  saving,
  selectedStock,
  onClose,
  onSubmit,
}: Omit<StockDetailsPanelProps, 'selectedStock'> & { selectedStock: InventoryStock }) {
  const [actionForm] = Form.useForm<CreateActionInput>();

  useEffect(() => {
    actionForm.setFieldsValue({ actionType: 'AWAITING_REVIEW', note: '' });
  }, [actionForm, selectedStock.id]);

  async function handleSubmit(values: CreateActionInput) {
    await onSubmit(values);
    actionForm.resetFields();
  }

  return (
    <aside className="h-full overflow-auto p-5">
      <div className="mb-4 flex items-start justify-between gap-3">
        <div>
          <Typography.Title level={4} className="mb-1!">
            {selectedStock.make} {selectedStock.model}
          </Typography.Title>
          <Typography.Text type="secondary">
            {selectedStock.modelYear} · {selectedStock.vin}
          </Typography.Text>
        </div>
        <Button onClick={onClose}>Close</Button>
      </div>

      <Space className="mb-4" wrap>
        <Tag color={selectedStock.status === 'IN_STOCK' ? 'green' : 'default'}>
          {selectedStock.status === 'IN_STOCK' ? 'In stock' : 'Out of stock'}
        </Tag>
        {selectedStock.isAging && <Tag color="red">Aging</Tag>}
        {selectedStock.latestAction && (
          <Tag color="blue">{actionLabels[selectedStock.latestAction.actionType]}</Tag>
        )}
      </Space>

      <Descriptions column={1} size="small" bordered>
        <Descriptions.Item label="Stock ID">{selectedStock.id}</Descriptions.Item>
        <Descriptions.Item label="Vehicle ID">{selectedStock.vehicleId}</Descriptions.Item>
        <Descriptions.Item label="Price">{currencyFormatter.format(selectedStock.price)}</Descriptions.Item>
        <Descriptions.Item label="Stock age">{selectedStock.inventoryAgeDays} days</Descriptions.Item>
        <Descriptions.Item label="Stocked in">
          {new Date(selectedStock.stockedInAt).toLocaleDateString()}
        </Descriptions.Item>
        {selectedStock.stockedOutAt && (
          <Descriptions.Item label="Stocked out">
            {new Date(selectedStock.stockedOutAt).toLocaleDateString()}
          </Descriptions.Item>
        )}
      </Descriptions>

      <Divider>Stock action</Divider>
      <Form
        form={actionForm}
        layout="vertical"
        initialValues={{ actionType: 'AWAITING_REVIEW', note: '' }}
        onFinish={handleSubmit}
      >
        <Form.Item name="actionType" label="Action" rules={[{ required: true }]}>
          <Select options={actionOptions} />
        </Form.Item>
        <Form.Item
          name="note"
          label="Note"
          rules={[{ required: true, message: 'Add a short note for the action.' }, { max: 500 }]}
        >
          <Input.TextArea rows={4} placeholder="Describe the next step and timing" />
        </Form.Item>
        <Button
          block
          type="primary"
          htmlType="submit"
          loading={saving}
          disabled={!selectedStock.isAging || selectedStock.status !== 'IN_STOCK'}
        >
          Save action
        </Button>
        {(!selectedStock.isAging || selectedStock.status !== 'IN_STOCK') && (
          <Typography.Text className="mt-2 block" type="secondary">
            Actions are available only for in-stock vehicles older than 90 days.
          </Typography.Text>
        )}
      </Form>

      <Divider>History</Divider>
      {history.length ? (
        <Timeline
          pending={historyLoading ? 'Loading history...' : undefined}
          items={history.map((event) => ({
            color: event.eventType === 'ACTION' ? 'blue' : event.eventType === 'STOCK_OUT' ? 'gray' : 'green',
            content: (
              <Space orientation="vertical" size={0}>
                <Typography.Text strong>
                  {event.eventType === 'ACTION' && event.actionType
                    ? actionLabels[event.actionType]
                    : event.eventType === 'STOCK_IN'
                      ? 'Stock in'
                      : 'Stock out'}
                </Typography.Text>
                <Typography.Text type="secondary">
                  {new Date(event.occurredAt).toLocaleString()}
                </Typography.Text>
                <Typography.Text>{event.note}</Typography.Text>
              </Space>
            ),
          }))}
        />
      ) : (
        <Empty description={historyLoading ? 'Loading history...' : 'No history yet'} />
      )}
    </aside>
  );
}
