# -*- coding: utf-8 -*-
"""
TencentBlueKing is pleased to support the open source community by making 蓝鲸智云-DB管理系统(BlueKing-BK-DBM) available.
Copyright (C) 2017-2023 THL A29 Limited, a Tencent company. All rights reserved.
Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
You may obtain a copy of the License at https://opensource.org/licenses/MIT
Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
"""
from django.utils.translation import ugettext_lazy as _
from rest_framework import serializers

from backend.db_meta.enums import ClusterPhase
from backend.flow.engine.controller.redis import RedisController
from backend.iam_app.dataclass.actions import ActionEnum
from backend.ticket import builders
from backend.ticket.builders.common.base import SkipToRepresentationMixin
from backend.ticket.builders.redis.base import BaseRedisTicketFlowBuilder, RedisSingleOpsBaseDetailSerializer
from backend.ticket.constants import TicketType


class RedisOpenDetailSerializer(RedisSingleOpsBaseDetailSerializer):
    pass


class RedisOpenFlowParamBuilder(builders.FlowParamBuilder):
    controller = RedisController.redis_cluster_open_close_scene

    def format_ticket_data(self):
        """
        {
            "uid": 340,
            "ticket_type": "PROXY_OPEN",
            "created_by": "admin",
            "cluster_id": 1111
        }
        """
        super().format_ticket_data()


@builders.BuilderFactory.register(
    TicketType.REDIS_PROXY_OPEN, phase=ClusterPhase.ONLINE, iam=ActionEnum.REDIS_OPEN_CLOSE
)
class RedisOpenFlowBuilder(BaseRedisTicketFlowBuilder):
    serializer = RedisOpenDetailSerializer
    inner_flow_builder = RedisOpenFlowParamBuilder
    inner_flow_name = _("启用集群")


class RedisInstanceOpenDetailSerializer(SkipToRepresentationMixin, serializers.Serializer):
    cluster_ids = serializers.ListField(help_text=_("集群ID列表"), child=serializers.IntegerField())
    force = serializers.BooleanField(help_text=_("是否强制"), required=False, default=True)


class RedisInstanceOpenFlowParamBuilder(builders.FlowParamBuilder):
    controller = RedisController.redis_ins_open_close_scene


@builders.BuilderFactory.register(
    TicketType.REDIS_INSTANCE_OPEN, phase=ClusterPhase.ONLINE, iam=ActionEnum.REDIS_OPEN_CLOSE
)
class RedisInstanceCloseFlowBuilder(BaseRedisTicketFlowBuilder):
    serializer = RedisInstanceOpenDetailSerializer
    inner_flow_builder = RedisInstanceOpenFlowParamBuilder
    inner_flow_name = _("启用集群")
