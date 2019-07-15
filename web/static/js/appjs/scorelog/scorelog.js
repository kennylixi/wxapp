var prefix = "/scorelog"
$(function() {
	load();
});

function load() {
	$('#exampleTable')
		.bootstrapTable(
			{
				method : 'get', // 服务器数据的请求方式 get or post
				url : prefix + "/list", // 服务器数据的加载地址
				// showRefresh : true,
				// showToggle : true,
				// showColumns : true,
				iconSize : 'outline',
				toolbar : '#exampleToolbar',
				striped : true, // 设置为true会有隔行变色效果
				dataType : "json", // 服务器返回的数据类型
				pagination : true, // 设置为true会在底部显示分页条
				// queryParamsType : "limit",
				// //设置为limit则会发送符合RESTFull格式的参数
				singleSelect : false, // 设置为true将禁止多选
				// contentType : "application/x-www-form-urlencoded",
				// //发送到服务器的数据编码类型
				pageSize : 10, // 如果设置了分页，每页数据条数
				pageNumber : 1, // 如果设置了分布，首页页码
				// search : true, // 是否显示搜索框
				showColumns : false, // 是否显示内容下拉框（选择显示的列）
				sidePagination : "server", // 设置在哪里进行分页，可选值为"client" 或者
				// "server"
				queryParams : function(params) {
					return {
						// 说明：传入后台的参数包括offset开始索引，limit步长，sort排序列，order：desc或者,以及所有列的键值对
						limit : params.limit,
						offset : params.offset,
						Account : $('#searchName').val(),
					};
				},
				// //请求服务器数据时，你可以通过重写参数的方式添加一些额外的参数，例如 toolbar 中的参数 如果
				// queryParamsType = 'limit' ,返回参数必须包含
				// limit, offset, search, sort, order 否则, 需要包含:
				// pageSize, pageNumber, searchText, sortName,
				// sortOrder.
				// 返回false将会终止请求
				columns : [
					{
						checkbox : true
					},
					{
						field : 'Id', // 列字段名
						title : '序号' // 列标题
					},
					{
						field : 'Name',
						title : '姓名'
					},
                    {
                        field : 'Account',
                        title : '账号'
                    },
                    {
                        field : 'Type',
                        title : '类型',
                        align : 'center',
                        formatter : function(value, row, index) {
                            if (value == '1') {
                                return '<span class="label label-primary">充值</span>';
                            } else if (value == '2') {
                                return '<span class="label label-danger">置顶</span>';
                            } else if (value == '3') {
                                return '<span class="label label-danger">发布</span>';
                            } else if (value == '4') {
                                return '<span class="label label-danger">更新</span>';
                            }
                        }
                    },
					{
						field : 'ChangeScore',
						title : '改变积分'
					},
					{
						field : 'Score',
						title : '总积分'
					},
                    {
                        field : 'Created',
                        title : '创建时间'
                    }]
			});
}
function reLoad() {
	$('#exampleTable').bootstrapTable('refresh');
}