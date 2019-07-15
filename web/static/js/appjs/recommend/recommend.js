var prefix = "/recommend"
var categoryId = '';
$(function() {
	load();
});

function load() {
	$('#exampleTable')
		.bootstrapTable(
			{
                cache : false,
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
                        Cid : categoryId,
                        Aid : $('#searchName').val(),
                        IsCost : $('#IsCost option:selected').val(),
                        Status : $('#Status option:selected').val(),
                        Type : $('#Type option:selected').val()
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
						field : 'Title',
						title : '标题'
					},
					{
                        field : 'Type',
                        title : '类型',
                        align : 'center',
                        formatter : function(value, row, index) {
                            if (value == '3') {
                                return '<span class="label label-warning">详情页推荐</span>';
                            } else if (value == '2') {
                                return '<span class="label label-info">首页分类推荐</span>';
                            } else if (value == '1') {
                                return '<span class="label label-success">首页推荐</span>';
                            }
                        }
					},
                    {
                        field : 'Status',
                        title : '是否上架',
                        align : 'center',
                        formatter : function(value, row, index) {
                            if (value == '3') {
                                return '<span class="label label-warning">已下架</span>';
                            } else if (value == '2') {
                                return '<span class="label label-info">已上架</span>';
                            } else if (value == '1') {
                                return '<span class="label label-success">未上架</span>';
                            }
                        }
                    },
                    {
                        field : 'IsCost',
                        title : '是否付费',
                        align : 'center',
                        formatter : function(value, row, index) {
                            if (value == '2') {
                                return '<span class="label label-info">是</span>';
                            } else if (value == '1') {
                                return '<span class="label label-success">否</span>';
                            }
                        }
                    },
                    {
                        field : 'Rtime',
                        align : 'center',
                        title : '推广时间(月)'
                    },
                    {
                        field : 'Etime',
                        title : '结束时间',
                        formatter: function (value, row, index) {
                            return changeDateFormat(row.Etime)
                        }
                    },
					{
						title : '操作',
						field : 'id',
						align : 'center',
						formatter : function(value, row, index) {
						    if (row.Status===1) {
                                var d = '<a class="btn btn-success btn-sm" href="#" title="上架"  mce_href="#" onclick="grounding(\''
                                    + row.Id
                                    + '\', \'上架\')"><i class="fa fa-arrow-up"></i></a> ';
                                return d;
                            }
                            if (row.Status===2) {
                                var d = '<a class="btn btn-success btn-sm" href="#" title="下架"  mce_href="#" onclick="grounding(\''
                                    + row.Id
                                    + '\', \'下架\')"><i class="fa fa-arrow-down"></i></a> ';
                                return d;
                            }
						}
					} ]
			});
}
function reLoad() {
	$('#exampleTable').bootstrapTable('refresh');
}
function add() {
	// iframe层
	layer.open({
		type : 2,
		title : '增加推荐',
		maxmin : true,
		shadeClose : false, // 点击遮罩关闭层
		area : [ '800px', '520px' ],
		content : prefix + '/add'
	});
}
function batchRemove() {
    var rows = $('#exampleTable').bootstrapTable('getSelections'); // 返回所有选择的行，当没有选择的记录时，返回一个空数组
    if (rows.length == 0) {
        layer.msg("请选择要删除的数据");
        return;
    }
    layer.confirm("确认要删除选中的'" + rows.length + "'条数据吗?", {
        btn : [ '确定', '取消' ]
        // 按钮
    }, function() {
        var ids = new Array();
        // 遍历所有选择的行数据，取每条数据对应的ID
        $.each(rows, function(i, row) {
            ids[i] = row['Id'];
        });
        $.ajax({
            type : 'POST',
            data : {
                "ids" : JSON.stringify(ids)
            },
            url : prefix + '/removes',
            success : function(r) {
                if (r.code == 0) {
                    layer.msg(r.msg);
                    reLoad();
                } else {
                    layer.msg(r.msg);
                }
            }
        });
    }, function() {});
}
function grounding(id, oper) {
    layer.confirm('确定要'+oper+'选中的记录？', {
        btn : [ '确定', '取消' ]
    }, function() {
        $.ajax({
            url : "/recommend/grounding",
            type : "post",
            data : {
                'id' : id
            },
            success : function(r) {
                if (r.code == 0) {
                    layer.msg(r.msg);
                    reLoad();
                } else {
                    layer.msg(r.msg);
                }
            }
        });
    })
}
var openCategory = function(){
    layer.open({
        type:2,
        title:"选择类目",
        area : [ '300px', '450px' ],
        content:"/category/treeview"
    })
};

function loadCategory( id, name){
    categoryId = id;
    $("#Cname").val(name);
    var opt = {
        query : {
            Cid : categoryId,
            Aid : $('#searchName').val(),
            IsCost : $('#IsCost option:selected').val(),
            Status : $('#Status option:selected').val(),
            Type : $('#Type option:selected').val()
        }
    }
    $('#exampleTable').bootstrapTable('refresh', opt);
}

$('#IsCost').on('change', function() {
    var opt = {
        query : {
            Cid : categoryId,
            Aid : $('#searchName').val(),
            IsCost : $('#IsCost option:selected').val(),
            Status : $('#Status option:selected').val(),
            Type : $('#Type option:selected').val()
        }
    }
    $('#exampleTable').bootstrapTable('refresh', opt);
});
$('#Status').on('change', function() {
    var opt = {
        query : {
            Cid : categoryId,
            Aid : $('#searchName').val(),
            IsCost : $('#IsCost option:selected').val(),
            Status : $('#Status option:selected').val(),
            Type : $('#Type option:selected').val()
        }
    }
    $('#exampleTable').bootstrapTable('refresh', opt);
});
$('#Type').on('change', function() {
    var type = $('#Type option:selected').val();
    if (type==="1"){
        categoryId = ""
        $("#Cname").val("所属类目");
    }
    var opt = {
        query : {
            Cid : categoryId,
            Aid : $('#searchName').val(),
            IsCost : $('#IsCost option:selected').val(),
            Status : $('#Status option:selected').val(),
            Type :type
        }
    }
    $('#exampleTable').bootstrapTable('refresh', opt);
});
//转换日期格式(时间戳转换为datetime格式)
function changeDateFormat(cellval) {
    if (cellval != null) {
        var date = new Date(cellval);
        if (date.getFullYear()===1){
            return ""
        }
        var month = date.getMonth() + 1 < 10 ? "0" + (date.getMonth() + 1) : date.getMonth() + 1;
        var currentDate = date.getDate() < 10 ? "0" + date.getDate() : date.getDate();

        var hours = date.getHours() < 10 ? "0" + date.getHours() : date.getHours();
        var minutes = date.getMinutes() < 10 ? "0" + date.getMinutes() : date.getMinutes();
        var seconds = date.getSeconds() < 10 ? "0" + date.getSeconds() : date.getSeconds();

        return date.getFullYear() + "-" + month + "-" + currentDate + " " + hours + ":" + minutes + ":" + seconds;
    }
}