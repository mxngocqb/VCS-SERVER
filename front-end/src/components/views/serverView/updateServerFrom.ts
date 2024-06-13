import { ServerStatus } from "@/plugins/axios/server/interfaces";
import { serverService } from "@/plugins/axios/server/serverService";
import { useMainStore } from "@/stores/mainStore";
import { useServerStore } from "@/stores/serverStore";
import { useForm } from "vee-validate";
import * as yup from "yup";

export default function () {
    const schema = yup.object({
        name: yup.string().required(),
        ipv4: yup
            .string()
            .matches(
                /^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$/,
                "This is not format of IPv4"
            )
            .required(),
        status: yup.boolean(),
    });

    const {
        handleSubmit,
        errors,
        isValidating,
        defineField,
        setFieldValue,
        resetForm,
    } = useForm({
        validationSchema: schema,
    });

    const [name, nameAttrs] = defineField("name");
    const [ipv4, ipv4Attrs] = defineField("ipv4");
    const mainStore = useMainStore();

    const onSubmit = handleSubmit(async (value) => {
        const serverStore = useServerStore();
        const selectedServer = serverStore.selectedServerComputed;
        if (selectedServer.value?.ID) {
            await serverService
                .updateServer({
                    id: selectedServer.value?.ID,
                    name: value["name"],
                    ipv4: value["ipv4"],
                })
                .then((res) => {
                    mainStore.showNofitication({
                        typeNotification: "infor",
                        title: "Update server",
                        content: "Update server succesfuly",
                    });
                });
        }
        return true;
    });
    return {
        resetForm,
        setFieldValue,
        name,
        ipv4,
        onSubmit,
        errors,
        isValidating,
    };
}
